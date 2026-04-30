package bo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/enums"
	"pix-psp-simulator/src/layers/main/models"
	"pix-psp-simulator/tests/testutil"
)

func TestCreateDevolucaoBO_Success(t *testing.T) {
	existingPix := &models.Pix{
		EndToEndID: "E123",
		TxID:       "txid123",
		Valor:      "100.00",
		Devolucoes: []models.Devolucao{},
	}
	pixRepo := &testutil.MockPixRepository{FindByE2EIDResult: existingPix}
	gen := &testutil.MockIDGenerator{RtrIDResult: "D123"}

	b := bo.NewCreateDevolucaoBO(pixRepo, gen)
	result, err := b.Execute(bo.CreateDevolucaoInput{
		E2EID: "E123",
		DevID: "dev01",
		Valor: "50.00",
	})

	require.NoError(t, err)
	assert.Equal(t, "dev01", result.Devolucao.ID)
	assert.Equal(t, "50.00", result.Devolucao.Valor)
	assert.Equal(t, enums.DevolucaoStatusDevolvido, result.Devolucao.Status)
}

func TestCreateDevolucaoBO_PaymentNotFound(t *testing.T) {
	pixRepo := &testutil.MockPixRepository{FindByE2EIDResult: nil}
	gen := &testutil.MockIDGenerator{}

	b := bo.NewCreateDevolucaoBO(pixRepo, gen)
	_, err := b.Execute(bo.CreateDevolucaoInput{
		E2EID: "naoexiste",
		DevID: "dev01",
		Valor: "50.00",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "pagamento não encontrado")
}

func TestCreateDevolucaoBO_ValueExceedsOriginal(t *testing.T) {
	existingPix := &models.Pix{
		EndToEndID: "E123",
		Valor:      "100.00",
		Devolucoes: []models.Devolucao{},
	}
	pixRepo := &testutil.MockPixRepository{FindByE2EIDResult: existingPix}
	gen := &testutil.MockIDGenerator{}

	b := bo.NewCreateDevolucaoBO(pixRepo, gen)
	_, err := b.Execute(bo.CreateDevolucaoInput{
		E2EID: "E123",
		DevID: "dev01",
		Valor: "150.00",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "excede o valor disponível")
}

func TestCreateDevolucaoBO_DuplicateDevID(t *testing.T) {
	existingPix := &models.Pix{
		EndToEndID: "E123",
		Valor:      "100.00",
		Devolucoes: []models.Devolucao{
			{ID: "dev01", Valor: "10.00", Status: enums.DevolucaoStatusDevolvido},
		},
	}
	pixRepo := &testutil.MockPixRepository{FindByE2EIDResult: existingPix}
	gen := &testutil.MockIDGenerator{}

	b := bo.NewCreateDevolucaoBO(pixRepo, gen)
	_, err := b.Execute(bo.CreateDevolucaoInput{
		E2EID: "E123",
		DevID: "dev01",
		Valor: "10.00",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "já existe")
}
