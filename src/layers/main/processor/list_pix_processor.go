package processor

import (
	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type ListPixRequest struct {
	TxID  string
	Inicio string
	Fim   string
}

type ListPixResponse struct {
	Pix []models.Pix
}

type ListPixProcessor struct {
	listPixBO *bo.ListPixBO
}

func NewListPixProcessor(listPixBO *bo.ListPixBO) *ListPixProcessor {
	return &ListPixProcessor{listPixBO: listPixBO}
}

func (p *ListPixProcessor) Process(req ListPixRequest) (*ListPixResponse, error) {
	output, err := p.listPixBO.Execute(bo.ListPixInput{
		TxID:  req.TxID,
		Inicio: req.Inicio,
		Fim:   req.Fim,
	})
	if err != nil {
		return nil, err
	}
	return &ListPixResponse{Pix: output.Pix}, nil
}
