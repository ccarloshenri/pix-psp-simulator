package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type GetDevolucaoRequest struct {
	E2EID string
	DevID string
}

type GetDevolucaoResponse struct {
	Devolucao models.Devolucao
}

type GetDevolucaoProcessor struct {
	getDevolucaoBO *bo.GetDevolucaoBO
}

func NewGetDevolucaoProcessor(getDevolucaoBO *bo.GetDevolucaoBO) *GetDevolucaoProcessor {
	return &GetDevolucaoProcessor{getDevolucaoBO: getDevolucaoBO}
}

func (p *GetDevolucaoProcessor) Process(req GetDevolucaoRequest) (*GetDevolucaoResponse, error) {
	if req.E2EID == "" {
		return nil, fmt.Errorf("e2eid é obrigatório")
	}
	if req.DevID == "" {
		return nil, fmt.Errorf("id da devolução é obrigatório")
	}

	output, err := p.getDevolucaoBO.Execute(bo.GetDevolucaoInput{
		E2EID: req.E2EID,
		DevID: req.DevID,
	})
	if err != nil {
		return nil, err
	}
	return &GetDevolucaoResponse{Devolucao: output.Devolucao}, nil
}
