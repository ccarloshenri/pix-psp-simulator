package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type CobValorRequest struct {
	Original string `json:"original"`
}

type CreateCobRequest struct {
	Chave          string                 `json:"chave"`
	Expiracao      int                    `json:"expiracao,omitempty"`
	Valor          CobValorRequest        `json:"valor"`
	Devedor        *models.Devedor        `json:"devedor,omitempty"`
	InfoAdicionais []models.InfoAdicional `json:"infoAdicionais,omitempty"`
}

type CreateCobResponse struct {
	Cob models.Cob
}

type CreateCobProcessor struct {
	createCobBO *bo.CreateCobBO
}

func NewCreateCobProcessor(createCobBO *bo.CreateCobBO) *CreateCobProcessor {
	return &CreateCobProcessor{createCobBO: createCobBO}
}

func (p *CreateCobProcessor) Process(txid string, req CreateCobRequest) (*CreateCobResponse, error) {
	if req.Chave == "" {
		return nil, fmt.Errorf("chave é obrigatória")
	}
	if req.Valor.Original == "" {
		return nil, fmt.Errorf("valor.original é obrigatório")
	}

	input := bo.CreateCobInput{
		TxID:           txid,
		Chave:          req.Chave,
		Expiracao:      req.Expiracao,
		Valor:          req.Valor.Original,
		Devedor:        req.Devedor,
		InfoAdicionais: req.InfoAdicionais,
	}

	output, err := p.createCobBO.Execute(input)
	if err != nil {
		return nil, err
	}
	return &CreateCobResponse{Cob: output.Cob}, nil
}
