package processor

import (
	"fmt"
	"regexp"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
)

var (
	regexValor = regexp.MustCompile(`^\d{1,10}\.\d{2}$`)
	regexCPF   = regexp.MustCompile(`^\d{11}$`)
	regexCNPJ  = regexp.MustCompile(`^\d{14}$`)
	regexTxID  = regexp.MustCompile(`^[a-zA-Z0-9]{26,35}$`)
)

type LocRequest struct {
	ID int `json:"id"`
}

type CreateCobCalendario struct {
	Expiracao int `json:"expiracao,omitempty"`
}

type CobValorRequest struct {
	Original            string           `json:"original"`
	ModalidadeAlteracao int              `json:"modalidadeAlteracao,omitempty"`
	Retirada            *models.Retirada `json:"retirada,omitempty"`
}

type CreateCobRequest struct {
	Calendario         *CreateCobCalendario   `json:"calendario,omitempty"`
	Devedor            *models.Devedor        `json:"devedor,omitempty"`
	Loc                *LocRequest            `json:"loc,omitempty"`
	Valor              CobValorRequest        `json:"valor"`
	Chave              string                 `json:"chave"`
	SolicitacaoPagador string                 `json:"solicitacaoPagador,omitempty"`
	InfoAdicionais     []models.InfoAdicional `json:"infoAdicionais,omitempty"`
}

type CreateCobResponse struct {
	Cob models.Cob
}

type CreateCobProcessor struct {
	createCobBO *bo.CreateCobBO
}

func NewCreateCobProcessor(createCobBO *bo.CreateCobBO) *CreateCobProcessor {
	return &CreateCobProcessor{createCobBO: createCobBO}
}

func (p *CreateCobProcessor) Process(txid string, req CreateCobRequest) (*CreateCobResponse, error) {
	if txid != "" && !regexTxID.MatchString(txid) {
		return nil, fmt.Errorf("txid deve conter entre 26 e 35 caracteres alfanuméricos [a-zA-Z0-9]")
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
	if req.Devedor != nil {
		if err := validateDevedor(req.Devedor.CPF, req.Devedor.CNPJ, req.Devedor.Nome); err != nil {
			return nil, err
		}
	}
	if len(req.SolicitacaoPagador) > 140 {
		return nil, fmt.Errorf("solicitacaoPagador deve ter no máximo 140 caracteres")
	}
	if err := validateInfoAdicionais(req.InfoAdicionais); err != nil {
		return nil, err
	}

	expiracao := 0
	if req.Calendario != nil {
		expiracao = req.Calendario.Expiracao
	}

	var locID int
	if req.Loc != nil {
		locID = req.Loc.ID
	}

	output, err := p.createCobBO.Execute(bo.CreateCobInput{
		TxID:                txid,
		Chave:               req.Chave,
		Expiracao:           expiracao,
		Valor:               req.Valor.Original,
		ModalidadeAlteracao: req.Valor.ModalidadeAlteracao,
		Retirada:            req.Valor.Retirada,
		Devedor:             req.Devedor,
		LocID:               locID,
		SolicitacaoPagador:  req.SolicitacaoPagador,
		InfoAdicionais:      req.InfoAdicionais,
	})
	if err != nil {
		return nil, err
	}
	return &CreateCobResponse{Cob: output.Cob}, nil
}

func validateDevedor(cpf, cnpj, nome string) error {
	if cpf != "" && cnpj != "" {
		return fmt.Errorf("devedor.cpf e devedor.cnpj não podem ser preenchidos ao mesmo tempo")
	}
	if cpf != "" && !regexCPF.MatchString(cpf) {
		return fmt.Errorf("devedor.cpf deve conter 11 dígitos numéricos")
	}
	if cnpj != "" && !regexCNPJ.MatchString(cnpj) {
		return fmt.Errorf("devedor.cnpj deve conter 14 dígitos numéricos")
	}
	if nome != "" && cpf == "" && cnpj == "" {
		return fmt.Errorf("devedor.nome preenchido requer devedor.cpf ou devedor.cnpj")
	}
	return nil
}

func validateInfoAdicionais(infos []models.InfoAdicional) error {
	if len(infos) > 50 {
		return fmt.Errorf("infoAdicionais deve ter no máximo 50 itens")
	}
	for _, info := range infos {
		if len(info.Nome) > 50 {
			return fmt.Errorf("infoAdicionais[].nome deve ter no máximo 50 caracteres")
		}
		if len(info.Valor) > 200 {
			return fmt.Errorf("infoAdicionais[].valor deve ter no máximo 200 caracteres")
		}
	}
	return nil
}
