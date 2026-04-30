package models

import "time"

type Pix struct {
	EndToEndID  string      `json:"endToEndId"`
	TxID        string      `json:"txid"`
	Valor       string      `json:"valor"`
	Horario     time.Time   `json:"horario"`
	Infopagador string      `json:"infoPagador,omitempty"`
	Devolucoes  []Devolucao `json:"devolucoes,omitempty"`
}
