package bo

import (
	"fmt"
	"time"

	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type SimulatePaymentInput struct {
	TxID        string
	Valor       string
	Infopagador string
}

type SimulatePaymentOutput struct {
	Pix models.Pix
}

type SimulatePaymentBO struct {
	cobRepo  interfaces.CobRepository
	cobvRepo interfaces.CobVRepository
	pixRepo  interfaces.PixRepository
	gen      interfaces.IDGenerator
}

func NewSimulatePaymentBO(
	cobRepo interfaces.CobRepository,
	cobvRepo interfaces.CobVRepository,
	pixRepo interfaces.PixRepository,
	gen interfaces.IDGenerator,
) *SimulatePaymentBO {
	return &SimulatePaymentBO{cobRepo: cobRepo, cobvRepo: cobvRepo, pixRepo: pixRepo, gen: gen}
}

func (b *SimulatePaymentBO) Execute(input SimulatePaymentInput) (*SimulatePaymentOutput, error) {
	cob, _ := b.cobRepo.FindByTxID(input.TxID)
	cobv, _ := b.cobvRepo.FindByTxID(input.TxID)

	if cob == nil && cobv == nil {
		return nil, fmt.Errorf("cobrança não encontrada para txid %s", input.TxID)
	}

	isCobV := cob == nil

	if !isCobV {
		if cob.Status != enums.CobStatusAtiva {
			return nil, fmt.Errorf("cobrança não está ativa")
		}
	} else {
		if cobv.Status != enums.CobStatusAtiva {
			return nil, fmt.Errorf("cobrança com vencimento não está ativa")
		}
	}

	e2eid := b.gen.GenerateE2EID("60746948")
	pix := models.Pix{
		EndToEndID:  e2eid,
		TxID:        input.TxID,
		Valor:       input.Valor,
		Horario:     time.Now().UTC(),
		Infopagador: input.Infopagador,
		Devolucoes:  []models.Devolucao{},
	}

	if err := b.pixRepo.Save(pix); err != nil {
		return nil, fmt.Errorf("erro ao registrar pagamento: %w", err)
	}

	if !isCobV {
		cob.Status = enums.CobStatusConcluida
		cob.Pix = append(cob.Pix, pix)
		b.cobRepo.Update(*cob)
	} else {
		cobv.Status = enums.CobStatusConcluida
		cobv.Pix = append(cobv.Pix, pix)
		b.cobvRepo.Update(*cobv)
	}

	return &SimulatePaymentOutput{Pix: pix}, nil
}
