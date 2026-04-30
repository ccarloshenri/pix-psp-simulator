package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type GetPixRequest struct {
	E2EID string
}

type GetPixResponse struct {
	Pix models.Pix
}

type GetPixProcessor struct {
	getPixBO *bo.GetPixBO
}

func NewGetPixProcessor(getPixBO *bo.GetPixBO) *GetPixProcessor {
	return &GetPixProcessor{getPixBO: getPixBO}
}

func (p *GetPixProcessor) Process(req GetPixRequest) (*GetPixResponse, error) {
	if req.E2EID == "" {
		return nil, fmt.Errorf("e2eid é obrigatório")
	}

	output, err := p.getPixBO.Execute(bo.GetPixInput{E2EID: req.E2EID})
	if err != nil {
		return nil, err
	}
	return &GetPixResponse{Pix: output.Pix}, nil
}
