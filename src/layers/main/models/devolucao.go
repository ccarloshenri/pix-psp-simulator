package models

import "time"

type Devolucao struct {
	ID                 string
	RtrID              string
	Valor              string
	Horario            HorarioDevolucao
	Status             string
	Natureza           string
	DescricaoDevolucao string
}

type HorarioDevolucao struct {
	Solicitacao time.Time
	Liquidacao  *time.Time
}
