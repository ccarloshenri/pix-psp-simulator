package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type ListCobInput struct {
	Status string
	Inicio string
	Fim    string
}

type ListCobOutput struct {
	Cobs []models.Cob
}

type ListCobBO struct {
	repo interfaces.CobRepository
}

func NewListCobBO(repo interfaces.CobRepository) *ListCobBO {
	return &ListCobBO{repo: repo}
}

func (b *ListCobBO) Execute(input ListCobInput) (*ListCobOutput, error) {
	cobs, err := b.repo.FindAll(interfaces.CobFilters{
		Status: input.Status,
		Inicio: input.Inicio,
		Fim:    input.Fim,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar cobranças: %w", err)
	}
	return &ListCobOutput{Cobs: cobs}, nil
}
