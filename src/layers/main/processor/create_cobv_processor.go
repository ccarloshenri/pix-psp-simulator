package processor

import (
	"fmt"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

type CreateCobVCalendario struct {
	DataDeVencimento       string `json:"dataDeVencimento"`
	ValidadeAposVencimento int    `json:"validadeAposVencimento,omitempty"`
}

type CobVValorRequest struct {
	Original   string                  `json:"original"`
	Multa      *models.ValorComponente `json:"multa,omitempty"`
	Juros      *models.ValorComponente `json:"juros,omitempty"`
	Abatimento *models.ValorComponente `json:"abatimento,omitempty"`
	Desconto   *models.DescontoCobV    `json:"desconto,omitempty"`
}

type CreateCobVRequest struct {
	Calendario         CreateCobVCalendario   `json:"calendario"`
	Devedor            models.Devedor         `json:"devedor"`
	Loc                *LocRequest            `json:"loc,omitempty"`
	Valor              CobVValorRequest       `json:"valor"`
	Chave              string                 `json:"chave"`
	SolicitacaoPagador string                 `json:"solicitacaoPagador,omitempty"`
	InfoAdicionais     []models.InfoAdicional `json:"infoAdicionais,omitempty"`
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
	if !regexTxID.MatchString(txid) {
		return nil, fmt.Errorf("txid deve conter entre 26 e 35 caracteres alfanuméricos [a-zA-Z0-9]")
	}
	if req.Calendario.DataDeVencimento == "" {
		return nil, fmt.Errorf("calendario.dataDeVencimento é obrigatório")
	}
	if err := validateDevedor(req.Devedor.CPF, req.Devedor.CNPJ, req.Devedor.Nome); err != nil {
		return nil, err
	}
	if req.Chave == "" {
		return nil, fmt.Errorf("chave é obrigatória")
	}
	if len(req.Chave) > 77 {
		return nil, fmt.Errorf("chave deve ter no máximo 77 caracteres")
	}
	if req.Valor.Original == "" {
		return nil, fmt.Errorf("valor.original é obrigatório")
	}
	if !regexValor.MatchString(req.Valor.Original) {
		return nil, fmt.Errorf("valor.original deve seguir o formato \\d{1,10}\\.\\d{2} (ex: \"100.00\")")
	}
	if len(req.SolicitacaoPagador) > 140 {
		return nil, fmt.Errorf("solicitacaoPagador deve ter no máximo 140 caracteres")
	}
	if err := validateInfoAdicionais(req.InfoAdicionais); err != nil {
		return nil, err
	}

	var locID int
	if req.Loc != nil {
		locID = req.Loc.ID
	}

	output, err := p.createCobVBO.Execute(bo.CreateCobVInput{
		TxID:                   txid,
		Chave:                  req.Chave,
		DataDeVencimento:       req.Calendario.DataDeVencimento,
		ValidadeAposVencimento: req.Calendario.ValidadeAposVencimento,
		Valor:                  req.Valor.Original,
		Devedor:                req.Devedor,
		LocID:                  locID,
		SolicitacaoPagador:     req.SolicitacaoPagador,
		InfoAdicionais:         req.InfoAdicionais,
		Multa:                  req.Valor.Multa,
		Juros:                  req.Valor.Juros,
		Abatimento:             req.Valor.Abatimento,
		Desconto:               req.Valor.Desconto,
	})
	if err != nil {
		return nil, err
	}
	return &CreateCobVResponse{CobV: output.CobV}, nil
}
