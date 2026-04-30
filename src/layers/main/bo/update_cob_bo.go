package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/interfaces"
)

type UpdateCobInput struct {
	TxID      string
	Valor     string
	Expiracao int
}

type UpdateCobOutput struct{}

type UpdateCobBO struct {
	repo interfaces.CobRepository
}

func NewUpdateCobBO(repo interfaces.CobRepository) *UpdateCobBO {
	return &UpdateCobBO{repo: repo}
}

func (b *UpdateCobBO) Execute(input UpdateCobInput) (*UpdateCobOutput, error) {
	cob, err := b.repo.FindByTxID(input.TxID)
	if err != nil || cob == nil {
		return nil, fmt.Errorf("cobrança não encontrada")
	}
	if cob.Status != enums.CobStatusAtiva {
		return nil, fmt.Errorf("só é possível alterar cobranças com status ATIVA")
	}
	if input.Valor != "" {
		cob.Valor.Original = input.Valor
	}
	if input.Expiracao > 0 {
		cob.Calendario.Expiracao = input.Expiracao
	}
	cob.Revisao++
	if err := b.repo.Update(*cob); err != nil {
		return nil, fmt.Errorf("erro ao atualizar cobrança: %w", err)
	}
	return &UpdateCobOutput{}, nil
}
