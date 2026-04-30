package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type CobVValorRequest struct {
	Original string                  `json:"original"`
	Multa    *models.ValorComponente `json:"multa,omitempty"`
	Juros    *models.ValorComponente `json:"juros,omitempty"`
	Desconto *models.ValorComponente `json:"desconto,omitempty"`
}

type CreateCobVRequest struct {
	Chave                  string                 `json:"chave"`
	DataDeVencimento       string                 `json:"dataDeVencimento"`
	ValidadeAposVencimento int                    `json:"validadeAposVencimento,omitempty"`
	Valor                  CobVValorRequest       `json:"valor"`
	Devedor                models.Devedor         `json:"devedor"`
	InfoAdicionais         []models.InfoAdicional `json:"infoAdicionais,omitempty"`
}

type CreateCobVResponse struct {
	CobV models.CobV
}

type CreateCobVProcessor struct {
	createCobVBO *bo.CreateCobVBO
}

func NewCreateCobVProcessor(createCobVBO *bo.CreateCobVBO) *CreateCobVProcessor {
	return &CreateCobVProcessor{createCobVBO: createCobVBO}
}

func (p *CreateCobVProcessor) Process(txid string, req CreateCobVRequest) (*CreateCobVResponse, error) {
	if txid == "" {
		return nil, fmt.Errorf("txid é obrigatório para cobranças com vencimento")
	}
	if req.Chave == "" {
		return nil, fmt.Errorf("chave é obrigatória")
	}
	if req.Valor.Original == "" {
		return nil, fmt.Errorf("valor.original é obrigatório")
	}
	if req.DataDeVencimento == "" {
		return nil, fmt.Errorf("dataDeVencimento é obrigatória")
	}

	output, err := p.createCobVBO.Execute(bo.CreateCobVInput{
		TxID:                   txid,
		Chave:                  req.Chave,
		DataDeVencimento:       req.DataDeVencimento,
		ValidadeAposVencimento: req.ValidadeAposVencimento,
		Valor:                  req.Valor.Original,
		Devedor:                req.Devedor,
		InfoAdicionais:         req.InfoAdicionais,
		Multa:                  req.Valor.Multa,
		Juros:                  req.Valor.Juros,
		Desconto:               req.Valor.Desconto,
	})
	if err != nil {
		return nil, err
	}
	return &CreateCobVResponse{CobV: output.CobV}, nil
}
