package models

import "time"

type CobV struct {
	TxID           string
	Calendario     CalendarioCobV
	Status         string
	Devedor        Devedor
	Valor          CobVValor
	Chave          string
	InfoAdicionais []InfoAdicional
	Location       string
	PixCopiaCola   string
	Pix            []Pix
}

type CalendarioCobV struct {
	DataDeVencimento       string // YYYY-MM-DD
	ValidadeAposVencimento int    // days
	Criacao                time.Time
}

type CobVValor struct {
	Original string
	Multa    *ValorComponente
	Juros    *ValorComponente
	Desconto *ValorComponente
}

type ValorComponente struct {
	Modalidade int
	ValorPerc  string
}
