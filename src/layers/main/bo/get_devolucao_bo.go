package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type GetDevolucaoInput struct {
	E2EID string
	DevID string
}

type GetDevolucaoOutput struct {
	Devolucao models.Devolucao
}

type GetDevolucaoBO struct {
	pixRepo interfaces.PixRepository
}

func NewGetDevolucaoBO(pixRepo interfaces.PixRepository) *GetDevolucaoBO {
	return &GetDevolucaoBO{pixRepo: pixRepo}
}

func (b *GetDevolucaoBO) Execute(input GetDevolucaoInput) (*GetDevolucaoOutput, error) {
	dev, err := b.pixRepo.FindDevolucao(input.E2EID, input.DevID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar devolução: %w", err)
	}
	if dev == nil {
		return nil, fmt.Errorf("devolução não encontrada")
	}
	return &GetDevolucaoOutput{Devolucao: *dev}, nil
}
