package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/interfaces"
)

type DeleteCobInput struct {
	TxID string
}

type DeleteCobOutput struct{}

type DeleteCobBO struct {
	repo interfaces.CobRepository
}

func NewDeleteCobBO(repo interfaces.CobRepository) *DeleteCobBO {
	return &DeleteCobBO{repo: repo}
}

func (b *DeleteCobBO) Execute(input DeleteCobInput) (*DeleteCobOutput, error) {
	cob, err := b.repo.FindByTxID(input.TxID)
	if err != nil || cob == nil {
		return nil, fmt.Errorf("cobrança não encontrada")
	}
	if cob.Status != enums.CobStatusAtiva {
		return nil, fmt.Errorf("só é possível remover cobranças com status ATIVA")
	}
	if err := b.repo.Delete(input.TxID); err != nil {
		return nil, fmt.Errorf("erro ao remover cobrança: %w", err)
	}
	return &DeleteCobOutput{}, nil
}
