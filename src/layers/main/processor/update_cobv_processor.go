package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
)

type UpdateCobVRequest struct {
	TxID                   string           `json:"-"`
	Valor                  CobVValorRequest `json:"valor,omitempty"`
	DataDeVencimento       string           `json:"dataDeVencimento,omitempty"`
	ValidadeAposVencimento int              `json:"validadeAposVencimento,omitempty"`
}

type UpdateCobVResponse struct{}

type UpdateCobVProcessor struct {
	updateCobVBO *bo.UpdateCobVBO
}

func NewUpdateCobVProcessor(updateCobVBO *bo.UpdateCobVBO) *UpdateCobVProcessor {
	return &UpdateCobVProcessor{updateCobVBO: updateCobVBO}
}

func (p *UpdateCobVProcessor) Process(req UpdateCobVRequest) (*UpdateCobVResponse, error) {
	if req.TxID == "" {
		return nil, fmt.Errorf("txid é obrigatório")
	}

	_, err := p.updateCobVBO.Execute(bo.UpdateCobVInput{
		TxID:                   req.TxID,
		Valor:                  req.Valor.Original,
		DataDeVencimento:       req.DataDeVencimento,
		ValidadeAposVencimento: req.ValidadeAposVencimento,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateCobVResponse{}, nil
}
