package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
)

type DeleteCobRequest struct {
	TxID string
}

type DeleteCobResponse struct{}

type DeleteCobProcessor struct {
	deleteCobBO *bo.DeleteCobBO
}

func NewDeleteCobProcessor(deleteCobBO *bo.DeleteCobBO) *DeleteCobProcessor {
	return &DeleteCobProcessor{deleteCobBO: deleteCobBO}
}

func (p *DeleteCobProcessor) Process(req DeleteCobRequest) (*DeleteCobResponse, error) {
	if req.TxID == "" {
		return nil, fmt.Errorf("txid é obrigatório")
	}

	_, err := p.deleteCobBO.Execute(bo.DeleteCobInput{TxID: req.TxID})
	if err != nil {
		return nil, err
	}
	return &DeleteCobResponse{}, nil
}
