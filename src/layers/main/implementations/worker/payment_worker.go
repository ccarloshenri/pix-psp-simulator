package worker

import (
	"context"
	"log/slog"
	"time"

	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

// PaymentWorker drains the PaymentQueue and persists each Pix payment,
// then marks the associated Cob or CobV as CONCLUIDA.
type PaymentWorker struct {
	queue    interfaces.PaymentQueue
	cobRepo  interfaces.CobRepository
	cobvRepo interfaces.CobVRepository
	pixRepo  interfaces.PixRepository
}

func NewPaymentWorker(
	queue interfaces.PaymentQueue,
	cobRepo interfaces.CobRepository,
	cobvRepo interfaces.CobVRepository,
	pixRepo interfaces.PixRepository,
) *PaymentWorker {
	return &PaymentWorker{
		queue:    queue,
		cobRepo:  cobRepo,
		cobvRepo: cobvRepo,
		pixRepo:  pixRepo,
	}
}

// Start blocks, processing jobs until ctx is cancelled. Run in a goroutine.
func (w *PaymentWorker) Start(ctx context.Context) {
	for {
		select {
		case job := <-w.queue.Jobs():
			w.process(job)
		case <-ctx.Done():
			return
		}
	}
}

func (w *PaymentWorker) process(job interfaces.PaymentJob) {
	pix := models.Pix{
		EndToEndID:  job.E2EID,
		TxID:        job.TxID,
		Valor:       job.Valor,
		Horario:     time.Now().UTC(),
		Infopagador: job.Infopagador,
		Devolucoes:  []models.Devolucao{},
	}

	if err := w.pixRepo.Save(pix); err != nil {
		slog.Error("failed to save pix", "e2eid", job.E2EID, "error", err)
		return
	}

	if cob, _ := w.cobRepo.FindByTxID(job.TxID); cob != nil {
		cob.Status = enums.CobStatusConcluida
		cob.Pix = append(cob.Pix, pix)
		if err := w.cobRepo.Update(*cob); err != nil {
			slog.Error("failed to update cob status", "txid", job.TxID, "error", err)
		}
		return
	}

	if cobv, _ := w.cobvRepo.FindByTxID(job.TxID); cobv != nil {
		cobv.Status = enums.CobStatusConcluida
		cobv.Pix = append(cobv.Pix, pix)
		if err := w.cobvRepo.Update(*cobv); err != nil {
			slog.Error("failed to update cobv status", "txid", job.TxID, "error", err)
		}
	}
}
