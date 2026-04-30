package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type GetPixInput struct {
	E2EID string
}

type GetPixOutput struct {
	Pix models.Pix
}

type GetPixBO struct {
	repo interfaces.PixRepository
}

func NewGetPixBO(repo interfaces.PixRepository) *GetPixBO {
	return &GetPixBO{repo: repo}
}

func (b *GetPixBO) Execute(input GetPixInput) (*GetPixOutput, error) {
	pix, err := b.repo.FindByE2EID(input.E2EID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pix: %w", err)
	}
	if pix == nil {
		return nil, fmt.Errorf("pix não encontrado")
	}
	return &GetPixOutput{Pix: *pix}, nil
}
