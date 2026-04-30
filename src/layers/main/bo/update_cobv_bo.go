package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/interfaces"
)

type UpdateCobVInput struct {
	TxID                   string
	Valor                  string
	DataDeVencimento       string
	ValidadeAposVencimento int
}

type UpdateCobVOutput struct{}

type UpdateCobVBO struct {
	repo interfaces.CobVRepository
}

func NewUpdateCobVBO(repo interfaces.CobVRepository) *UpdateCobVBO {
	return &UpdateCobVBO{repo: repo}
}

func (b *UpdateCobVBO) Execute(input UpdateCobVInput) (*UpdateCobVOutput, error) {
	cobv, err := b.repo.FindByTxID(input.TxID)
	if err != nil || cobv == nil {
		return nil, fmt.Errorf("cobrança com vencimento não encontrada")
	}
	if cobv.Status != enums.CobStatusAtiva {
		return nil, fmt.Errorf("só é possível alterar cobranças com vencimento com status ATIVA")
	}
	if input.Valor != "" {
		cobv.Valor.Original = input.Valor
	}
	if input.DataDeVencimento != "" {
		cobv.Calendario.DataDeVencimento = input.DataDeVencimento
	}
	if input.ValidadeAposVencimento > 0 {
		cobv.Calendario.ValidadeAposVencimento = input.ValidadeAposVencimento
	}
	if err := b.repo.Update(*cobv); err != nil {
		return nil, fmt.Errorf("erro ao atualizar cobrança com vencimento: %w", err)
	}
	return &UpdateCobVOutput{}, nil
}
