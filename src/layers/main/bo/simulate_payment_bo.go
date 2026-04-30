package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/interfaces"
)

type SimulatePaymentInput struct {
	TxID        string
	Valor       string
	Infopagador string
}

type SimulatePaymentOutput struct {
	EndToEndID string
	TxID       string
}

type SimulatePaymentBO struct {
	cobRepo  interfaces.CobRepository
	cobvRepo interfaces.CobVRepository
	gen      interfaces.IDGenerator
	queue    interfaces.PaymentQueue
}

func NewSimulatePaymentBO(
	cobRepo interfaces.CobRepository,
	cobvRepo interfaces.CobVRepository,
	gen interfaces.IDGenerator,
	queue interfaces.PaymentQueue,
) *SimulatePaymentBO {
	return &SimulatePaymentBO{cobRepo: cobRepo, cobvRepo: cobvRepo, gen: gen, queue: queue}
}

func (b *SimulatePaymentBO) Execute(input SimulatePaymentInput) (*SimulatePaymentOutput, error) {
	cob, _ := b.cobRepo.FindByTxID(input.TxID)
	cobv, _ := b.cobvRepo.FindByTxID(input.TxID)

	if cob == nil && cobv == nil {
		return nil, fmt.Errorf("cobrança não encontrada para txid %s", input.TxID)
	}

	if cob != nil && cob.Status != enums.CobStatusAtiva {
		return nil, fmt.Errorf("cobrança não está ativa")
	}

	if cobv != nil && cob == nil && cobv.Status != enums.CobStatusAtiva {
		return nil, fmt.Errorf("cobrança com vencimento não está ativa")
	}

	e2eid := b.gen.GenerateE2EID("60746948")
	b.queue.Enqueue(interfaces.PaymentJob{
		E2EID:       e2eid,
		TxID:        input.TxID,
		Valor:       input.Valor,
		Infopagador: input.Infopagador,
	})

	return &SimulatePaymentOutput{EndToEndID: e2eid, TxID: input.TxID}, nil
}
