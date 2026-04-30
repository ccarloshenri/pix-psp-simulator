package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
)

type UpdateCobRequest struct {
	TxID      string
	Valor     CobValorRequest `json:"valor,omitempty"`
	Expiracao int             `json:"expiracao,omitempty"`
}

type UpdateCobResponse struct{}

type UpdateCobProcessor struct {
	updateCobBO *bo.UpdateCobBO
}

func NewUpdateCobProcessor(updateCobBO *bo.UpdateCobBO) *UpdateCobProcessor {
	return &UpdateCobProcessor{updateCobBO: updateCobBO}
}

func (p *UpdateCobProcessor) Process(req UpdateCobRequest) (*UpdateCobResponse, error) {
	if req.TxID == "" {
		return nil, fmt.Errorf("txid é obrigatório")
	}

	_, err := p.updateCobBO.Execute(bo.UpdateCobInput{
		TxID:      req.TxID,
		Valor:     req.Valor.Original,
		Expiracao: req.Expiracao,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateCobResponse{}, nil
}
