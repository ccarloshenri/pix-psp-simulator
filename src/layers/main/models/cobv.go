package models

import "time"

type CobV struct {
	TxID               string          `json:"txid"`
	Calendario         CalendarioCobV  `json:"calendario"`
	Revisao            int             `json:"revisao"`
	Status             string          `json:"status"`
	Devedor            Devedor         `json:"devedor"`
	Loc                *Loc            `json:"loc,omitempty"`
	Valor              CobVValor       `json:"valor"`
	PixCopiaECola      string          `json:"pixCopiaECola,omitempty"`
	Chave              string          `json:"chave"`
	SolicitacaoPagador string          `json:"solicitacaoPagador,omitempty"`
	InfoAdicionais     []InfoAdicional `json:"infoAdicionais,omitempty"`
	Pix                []Pix           `json:"pix,omitempty"`
}

type CalendarioCobV struct {
	DataDeVencimento       string    `json:"dataDeVencimento"`
	ValidadeAposVencimento int       `json:"validadeAposVencimento,omitempty"`
	Criacao                time.Time `json:"criacao"`
}

type CobVValor struct {
	Original   string           `json:"original"`
	Multa      *ValorComponente `json:"multa,omitempty"`
	Juros      *ValorComponente `json:"juros,omitempty"`
	Abatimento *ValorComponente `json:"abatimento,omitempty"`
	Desconto   *DescontoCobV    `json:"desconto,omitempty"`
}

type ValorComponente struct {
	Modalidade int    `json:"modalidade"`
	ValorPerc  string `json:"valorPerc"`
}

type DescontoCobV struct {
	Modalidade       int                `json:"modalidade"`
	ValorPerc        string             `json:"valorPerc,omitempty"`
	DescontoDataFixa []DescontoDataFixa `json:"descontoDataFixa,omitempty"`
}

type DescontoDataFixa struct {
	Data      string `json:"data"`
	ValorPerc string `json:"valorPerc"`
}
