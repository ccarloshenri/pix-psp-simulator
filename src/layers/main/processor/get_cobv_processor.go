package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type GetCobVRequest struct {
	TxID string
}

type GetCobVResponse struct {
	CobV models.CobV
}

type GetCobVProcessor struct {
	getCobVBO *bo.GetCobVBO
}

func NewGetCobVProcessor(getCobVBO *bo.GetCobVBO) *GetCobVProcessor {
	return &GetCobVProcessor{getCobVBO: getCobVBO}
}

func (p *GetCobVProcessor) Process(req GetCobVRequest) (*GetCobVResponse, error) {
	if req.TxID == "" {
		return nil, fmt.Errorf("txid é obrigatório")
	}

	output, err := p.getCobVBO.Execute(bo.GetCobVInput{TxID: req.TxID})
	if err != nil {
		return nil, err
	}
	return &GetCobVResponse{CobV: output.CobV}, nil
}
