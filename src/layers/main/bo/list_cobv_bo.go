package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type ListCobVInput struct {
	Status           string
	DataDeVencimento string
	Inicio           string
	Fim              string
}

type ListCobVOutput struct {
	CobVs []models.CobV
}

type ListCobVBO struct {
	repo interfaces.CobVRepository
}

func NewListCobVBO(repo interfaces.CobVRepository) *ListCobVBO {
	return &ListCobVBO{repo: repo}
}

func (b *ListCobVBO) Execute(input ListCobVInput) (*ListCobVOutput, error) {
	cobvs, err := b.repo.FindAll(interfaces.CobVFilters{
		Status:           input.Status,
		DataDeVencimento: input.DataDeVencimento,
		Inicio:           input.Inicio,
		Fim:              input.Fim,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar cobranças com vencimento: %w", err)
	}
	return &ListCobVOutput{CobVs: cobvs}, nil
}
