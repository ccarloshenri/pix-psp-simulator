package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type GetCobVInput struct {
	TxID string
}

type GetCobVOutput struct {
	CobV models.CobV
}

type GetCobVBO struct {
	repo interfaces.CobVRepository
}

func NewGetCobVBO(repo interfaces.CobVRepository) *GetCobVBO {
	return &GetCobVBO{repo: repo}
}

func (b *GetCobVBO) Execute(input GetCobVInput) (*GetCobVOutput, error) {
	cobv, err := b.repo.FindByTxID(input.TxID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cobrança com vencimento: %w", err)
	}
	if cobv == nil {
		return nil, fmt.Errorf("cobrança com vencimento não encontrada")
	}
	return &GetCobVOutput{CobV: *cobv}, nil
}
