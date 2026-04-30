package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type GetCobRequest struct {
	TxID string
}

type GetCobResponse struct {
	Cob models.Cob
}

type GetCobProcessor struct {
	getCobBO *bo.GetCobBO
}

func NewGetCobProcessor(getCobBO *bo.GetCobBO) *GetCobProcessor {
	return &GetCobProcessor{getCobBO: getCobBO}
}

func (p *GetCobProcessor) Process(req GetCobRequest) (*GetCobResponse, error) {
	if req.TxID == "" {
		return nil, fmt.Errorf("txid é obrigatório")
	}

	output, err := p.getCobBO.Execute(bo.GetCobInput{TxID: req.TxID})
	if err != nil {
		return nil, err
	}
	return &GetCobResponse{Cob: output.Cob}, nil
}
