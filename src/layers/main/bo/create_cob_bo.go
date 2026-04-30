package bo

import (
	"fmt"
	"time"

	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type CreateCobInput struct {
	TxID                string
	Chave               string
	Expiracao           int
	Valor               string
	ModalidadeAlteracao int
	Retirada            *models.Retirada
	Devedor             *models.Devedor
	LocID               int
	SolicitacaoPagador  string
	InfoAdicionais      []models.InfoAdicional
}

type CreateCobOutput struct {
	Cob models.Cob
}

type CreateCobBO struct {
	repo interfaces.CobRepository
	gen  interfaces.IDGenerator
}

func NewCreateCobBO(repo interfaces.CobRepository, gen interfaces.IDGenerator) *CreateCobBO {
	return &CreateCobBO{repo: repo, gen: gen}
}

func (b *CreateCobBO) Execute(input CreateCobInput) (*CreateCobOutput, error) {
	txid := input.TxID
	if txid == "" {
		txid = b.gen.GenerateTxID()
	}

	existing, _ := b.repo.FindByTxID(txid)
	if existing != nil {
		return nil, fmt.Errorf("cobrança com txid %s já existe", txid)
	}

	expiracao := input.Expiracao
	if expiracao <= 0 {
		expiracao = 86400
	}

	now := time.Now().UTC()

	var loc *models.Loc
	if input.LocID > 0 {
		loc = &models.Loc{
			ID:       input.LocID,
			Location: fmt.Sprintf("pix.simulator/cobqrcode/%s", txid),
			TipoCob:  "cob",
			Criacao:  now,
		}
	}

	cob := models.Cob{
		TxID: txid,
		Calendario: models.Calendario{
			Criacao:   now,
			Expiracao: expiracao,
		},
		Revisao: 0,
		Status:  enums.CobStatusAtiva,
		Chave:   input.Chave,
		Devedor: input.Devedor,
		Loc:     loc,
		Location: fmt.Sprintf("pix.simulator/cobqrcode/%s", txid),
		Valor: models.CobValor{
			Original:            input.Valor,
			ModalidadeAlteracao: input.ModalidadeAlteracao,
			Retirada:            input.Retirada,
		},
		PixCopiaECola:      fmt.Sprintf("00020126580014br.gov.bcb.pix0136%s5204000053039865802BR5925Simulador PSP PIX6009SAO PAULO62290525%s6304", input.Chave, txid),
		SolicitacaoPagador: input.SolicitacaoPagador,
		InfoAdicionais:     input.InfoAdicionais,
		Pix:                []models.Pix{},
	}

	if err := b.repo.Save(cob); err != nil {
		return nil, fmt.Errorf("erro ao salvar cobrança: %w", err)
	}

	return &CreateCobOutput{Cob: cob}, nil
}
