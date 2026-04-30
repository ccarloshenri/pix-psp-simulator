package models

import "time"

type Devolucao struct {
	ID                 string           `json:"id"`
	RtrID              string           `json:"rtrId"`
	Valor              string           `json:"valor"`
	Horario            HorarioDevolucao `json:"horario"`
	Status             string           `json:"status"`
	Natureza           string           `json:"natureza,omitempty"`
	DescricaoDevolucao string           `json:"descricaoDevolucao,omitempty"`
}

type HorarioDevolucao struct {
	Solicitacao time.Time  `json:"solicitacao"`
	Liquidacao  *time.Time `json:"liquidacao,omitempty"`
}
