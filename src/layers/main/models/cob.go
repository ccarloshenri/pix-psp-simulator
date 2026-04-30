package models

import "time"

type Cob struct {
	TxID               string          `json:"txid"`
	Calendario         Calendario      `json:"calendario"`
	Revisao            int             `json:"revisao"`
	Status             string          `json:"status"`
	Devedor            *Devedor        `json:"devedor,omitempty"`
	Loc                *Loc            `json:"loc,omitempty"`
	Location           string          `json:"location,omitempty"`
	Valor              CobValor        `json:"valor"`
	PixCopiaECola      string          `json:"pixCopiaECola,omitempty"`
	Chave              string          `json:"chave"`
	SolicitacaoPagador string          `json:"solicitacaoPagador,omitempty"`
	InfoAdicionais     []InfoAdicional `json:"infoAdicionais,omitempty"`
	Pix                []Pix           `json:"pix,omitempty"`
}

type Calendario struct {
	Criacao   time.Time `json:"criacao"`
	Expiracao int       `json:"expiracao"`
}

type CobValor struct {
	Original            string    `json:"original"`
	ModalidadeAlteracao int       `json:"modalidadeAlteracao,omitempty"`
	Retirada            *Retirada `json:"retirada,omitempty"`
}

type Loc struct {
	ID       int       `json:"id"`
	Location string    `json:"location,omitempty"`
	TipoCob  string    `json:"tipoCob,omitempty"`
	Criacao  time.Time `json:"criacao,omitempty"`
}

type Retirada struct {
	Saque *Saque `json:"saque,omitempty"`
	Troco *Troco `json:"troco,omitempty"`
}

type Saque struct {
	Valor                     string `json:"valor"`
	ModalidadeAlteracao       int    `json:"modalidadeAlteracao,omitempty"`
	ModalidadeAgente          string `json:"modalidadeAgente"`
	PrestadorDoServicoDeSaque string `json:"prestadorDoServicoDeSaque"`
}

type Troco struct {
	Valor                     string `json:"valor"`
	ModalidadeAlteracao       int    `json:"modalidadeAlteracao,omitempty"`
	ModalidadeAgente          string `json:"modalidadeAgente"`
	PrestadorDoServicoDeSaque string `json:"prestadorDoServicoDeSaque"`
}
