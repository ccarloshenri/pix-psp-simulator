package models

import "time"

type Pix struct {
	EndToEndID  string
	TxID        string
	Valor       string
	Horario     time.Time
	Infopagador string
	Devolucoes  []Devolucao
}
