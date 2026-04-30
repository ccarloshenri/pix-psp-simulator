package bo

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type ListPixInput struct {
	TxID  string
	Inicio string
	Fim   string
}

type ListPixOutput struct {
	Pix []models.Pix
}

type ListPixBO struct {
	repo interfaces.PixRepository
}

func NewListPixBO(repo interfaces.PixRepository) *ListPixBO {
	return &ListPixBO{repo: repo}
}

func (b *ListPixBO) Execute(input ListPixInput) (*ListPixOutput, error) {
	pix, err := b.repo.FindAll(interfaces.PixFilters{
		TxID:  input.TxID,
		Inicio: input.Inicio,
		Fim:   input.Fim,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pagamentos: %w", err)
	}
	return &ListPixOutput{Pix: pix}, nil
}
