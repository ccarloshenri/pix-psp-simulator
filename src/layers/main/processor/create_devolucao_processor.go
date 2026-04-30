package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type CreateDevolucaoRequest struct {
	Valor              string `json:"valor"`
	Natureza           string `json:"natureza,omitempty"`
	DescricaoDevolucao string `json:"descricaoDevolucao,omitempty"`
}

type CreateDevolucaoResponse struct {
	Devolucao models.Devolucao
}

type CreateDevolucaoProcessor struct {
	createDevolucaoBO *bo.CreateDevolucaoBO
}

func NewCreateDevolucaoProcessor(createDevolucaoBO *bo.CreateDevolucaoBO) *CreateDevolucaoProcessor {
	return &CreateDevolucaoProcessor{createDevolucaoBO: createDevolucaoBO}
}

func (p *CreateDevolucaoProcessor) Process(e2eid, devID string, req CreateDevolucaoRequest) (*CreateDevolucaoResponse, error) {
	if e2eid == "" {
		return nil, fmt.Errorf("e2eid é obrigatório")
	}
	if devID == "" {
		return nil, fmt.Errorf("id da devolução é obrigatório")
	}
	if req.Valor == "" {
		return nil, fmt.Errorf("valor é obrigatório")
	}

	output, err := p.createDevolucaoBO.Execute(bo.CreateDevolucaoInput{
		E2EID:              e2eid,
		DevID:              devID,
		Valor:              req.Valor,
		Natureza:           req.Natureza,
		DescricaoDevolucao: req.DescricaoDevolucao,
	})
	if err != nil {
		return nil, err
	}
	return &CreateDevolucaoResponse{Devolucao: output.Devolucao}, nil
}
