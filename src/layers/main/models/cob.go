package models

import "time"

type Cob struct {
	TxID           string
	Calendario     Calendario
	Status         string
	Devedor        *Devedor
	Valor          CobValor
	Chave          string
	InfoAdicionais []InfoAdicional
	Location       string
	PixCopiaCola   string
	Pix            []Pix
}

type Calendario struct {
	Criacao   time.Time
	Expiracao int // seconds, default 3600
}

type CobValor struct {
	Original   string
	Modalidade int // 0=fixed
}
