package bo

import (
	"fmt"
	"time"

	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type CreateCobVInput struct {
	TxID                   string
	Chave                  string
	DataDeVencimento       string
	ValidadeAposVencimento int
	Valor                  string
	Devedor                models.Devedor
	InfoAdicionais         []models.InfoAdicional
	Multa                  *models.ValorComponente
	Juros                  *models.ValorComponente
	Desconto               *models.ValorComponente
}

type CreateCobVOutput struct {
	CobV models.CobV
}

type CreateCobVBO struct {
	repo interfaces.CobVRepository
	gen  interfaces.IDGenerator
}

func NewCreateCobVBO(repo interfaces.CobVRepository, gen interfaces.IDGenerator) *CreateCobVBO {
	return &CreateCobVBO{repo: repo, gen: gen}
}

func (b *CreateCobVBO) Execute(input CreateCobVInput) (*CreateCobVOutput, error) {
	existing, _ := b.repo.FindByTxID(input.TxID)
	if existing != nil {
		return nil, fmt.Errorf("cobrança com vencimento com txid %s já existe", input.TxID)
	}

	validadeApos := input.ValidadeAposVencimento
	if validadeApos <= 0 {
		validadeApos = 30
	}

	cobv := models.CobV{
		TxID: input.TxID,
		Calendario: models.CalendarioCobV{
			DataDeVencimento:       input.DataDeVencimento,
			ValidadeAposVencimento: validadeApos,
			Criacao:                time.Now().UTC(),
		},
		Status:  enums.CobStatusAtiva,
		Devedor: input.Devedor,
		Valor: models.CobVValor{
			Original: input.Valor,
			Multa:    input.Multa,
			Juros:    input.Juros,
			Desconto: input.Desconto,
		},
		Chave:          input.Chave,
		InfoAdicionais: input.InfoAdicionais,
		Location:       fmt.Sprintf("pix.simulator/cobvqrcode/%s", input.TxID),
		PixCopiaCola:   fmt.Sprintf("00020126580014br.gov.bcb.pix0136%s5204000053039865802BR5925Simulador PSP PIX6009SAO PAULO62290525%s6304", input.Chave, input.TxID),
		Pix:            []models.Pix{},
	}

	if err := b.repo.Save(cobv); err != nil {
		return nil, fmt.Errorf("erro ao salvar cobrança com vencimento: %w", err)
	}

	return &CreateCobVOutput{CobV: cobv}, nil
}
