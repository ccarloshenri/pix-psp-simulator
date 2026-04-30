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

func TestCreateCobBO_Success(t *testing.T) {
	repo := &testutil.MockCobRepository{}
	gen := &testutil.MockIDGenerator{TxIDResult: "txid123"}

	b := bo.NewCreateCobBO(repo, gen)
	result, err := b.Execute(bo.CreateCobInput{
		Chave: "+5511999998888",
		Valor: "100.00",
	})

	require.NoError(t, err)
	assert.Equal(t, "txid123", result.Cob.TxID)
	assert.Equal(t, enums.CobStatusAtiva, result.Cob.Status)
	assert.Equal(t, "100.00", result.Cob.Valor.Original)
	assert.True(t, repo.SaveCalled)
}

func TestCreateCobBO_UseProvidedTxID(t *testing.T) {
	repo := &testutil.MockCobRepository{}
	gen := &testutil.MockIDGenerator{TxIDResult: "generated"}

	b := bo.NewCreateCobBO(repo, gen)
	result, err := b.Execute(bo.CreateCobInput{
		TxID:  "meutxid",
		Chave: "+5511999998888",
		Valor: "50.00",
	})

	require.NoError(t, err)
	assert.Equal(t, "meutxid", result.Cob.TxID)
}

func TestCreateCobBO_DefaultExpiracao(t *testing.T) {
	repo := &testutil.MockCobRepository{}
	gen := &testutil.MockIDGenerator{TxIDResult: "txid123"}

	b := bo.NewCreateCobBO(repo, gen)
	result, err := b.Execute(bo.CreateCobInput{
		Chave: "+5511999998888",
		Valor: "100.00",
	})

	require.NoError(t, err)
	assert.Equal(t, 86400, result.Cob.Calendario.Expiracao)
}

func TestCreateCobBO_DuplicateTxIDReturnsError(t *testing.T) {
	existing := &models.Cob{TxID: "txid123"}
	repo := &testutil.MockCobRepository{FindResult: existing}
	gen := &testutil.MockIDGenerator{TxIDResult: "txid123"}

	b := bo.NewCreateCobBO(repo, gen)
	_, err := b.Execute(bo.CreateCobInput{
		TxID:  "txid123",
		Chave: "+5511999998888",
		Valor: "100.00",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "já existe")
}

func TestCreateCobBO_SaveError(t *testing.T) {
	repo := &testutil.MockCobRepository{SaveError: testutil.ErrGeneric}
	gen := &testutil.MockIDGenerator{TxIDResult: "txid123"}

	b := bo.NewCreateCobBO(repo, gen)
	_, err := b.Execute(bo.CreateCobInput{
		Chave: "+5511999998888",
		Valor: "100.00",
	})

	require.ErrorIs(t, err, testutil.ErrGeneric)
}
