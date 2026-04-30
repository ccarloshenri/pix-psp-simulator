package models

type Devedor struct {
	CPF        string `json:"cpf,omitempty"`
	CNPJ       string `json:"cnpj,omitempty"`
	Nome       string `json:"nome,omitempty"`
	Email      string `json:"email,omitempty"`
	Logradouro string `json:"logradouro,omitempty"`
	Cidade     string `json:"cidade,omitempty"`
	UF         string `json:"uf,omitempty"`
	CEP        string `json:"cep,omitempty"`
}

type InfoAdicional struct {
	Nome  string `json:"nome"`
	Valor string `json:"valor"`
}
