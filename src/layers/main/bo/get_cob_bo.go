package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type GetCobInput struct {
	TxID string
}

type GetCobOutput struct {
	Cob models.Cob
}

type GetCobBO struct {
	repo interfaces.CobRepository
}

func NewGetCobBO(repo interfaces.CobRepository) *GetCobBO {
	return &GetCobBO{repo: repo}
}

func (b *GetCobBO) Execute(input GetCobInput) (*GetCobOutput, error) {
	cob, err := b.repo.FindByTxID(input.TxID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cobrança: %w", err)
	}
	if cob == nil {
		return nil, fmt.Errorf("cobrança não encontrada")
	}
	return &GetCobOutput{Cob: *cob}, nil
}
