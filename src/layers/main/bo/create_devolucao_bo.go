package bo

import (
	"fmt"
	"strconv"
	"time"

	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type CreateDevolucaoInput struct {
	E2EID              string
	DevID              string
	Valor              string
	Natureza           string
	DescricaoDevolucao string
}

type CreateDevolucaoOutput struct {
	Devolucao models.Devolucao
}

type CreateDevolucaoBO struct {
	pixRepo interfaces.PixRepository
	gen     interfaces.IDGenerator
}

func NewCreateDevolucaoBO(pixRepo interfaces.PixRepository, gen interfaces.IDGenerator) *CreateDevolucaoBO {
	return &CreateDevolucaoBO{pixRepo: pixRepo, gen: gen}
}

func (b *CreateDevolucaoBO) Execute(input CreateDevolucaoInput) (*CreateDevolucaoOutput, error) {
	pix, err := b.pixRepo.FindByE2EID(input.E2EID)
	if err != nil || pix == nil {
		return nil, fmt.Errorf("pagamento não encontrado")
	}

	for _, d := range pix.Devolucoes {
		if d.ID == input.DevID {
			return nil, fmt.Errorf("devolução com id %s já existe", input.DevID)
		}
	}

	originalVal, err := strconv.ParseFloat(pix.Valor, 64)
	if err != nil {
		return nil, fmt.Errorf("valor original inválido")
	}
	refundVal, err := strconv.ParseFloat(input.Valor, 64)
	if err != nil {
		return nil, fmt.Errorf("valor de devolução inválido")
	}

	var totalRefunded float64
	for _, d := range pix.Devolucoes {
		if d.Status == enums.DevolucaoStatusDevolvido || d.Status == enums.DevolucaoStatusEmProcessamento {
			v, _ := strconv.ParseFloat(d.Valor, 64)
			totalRefunded += v
		}
	}

	if refundVal+totalRefunded > originalVal {
		return nil, fmt.Errorf("valor de devolução excede o valor disponível para devolução")
	}

	natureza := input.Natureza
	if natureza == "" {
		natureza = "ORIGINAL"
	}

	now := time.Now().UTC()
	rtrID := b.gen.GenerateRtrID("60746948")

	dev := models.Devolucao{
		ID:    input.DevID,
		RtrID: rtrID,
		Valor: input.Valor,
		Horario: models.HorarioDevolucao{
			Solicitacao: now,
		},
		Status:             enums.DevolucaoStatusEmProcessamento,
		Natureza:           natureza,
		DescricaoDevolucao: input.DescricaoDevolucao,
	}

	if err := b.pixRepo.AddDevolucao(input.E2EID, dev); err != nil {
		return nil, fmt.Errorf("erro ao registrar devolução: %w", err)
	}

	// Simulate immediate settlement — mark as DEVOLVIDO
	liquidacao := time.Now().UTC().Add(2 * time.Second)
	dev.Status = enums.DevolucaoStatusDevolvido
	dev.Horario.Liquidacao = &liquidacao
	b.pixRepo.UpdateDevolucao(input.E2EID, dev)

	return &CreateDevolucaoOutput{Devolucao: dev}, nil
}
