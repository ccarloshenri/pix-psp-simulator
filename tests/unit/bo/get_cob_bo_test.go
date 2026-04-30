package bo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
	"pix-psp-simulator/tests/testutil"
)

func TestGetCobBO_Found(t *testing.T) {
	expected := &models.Cob{TxID: "txid123", Chave: "+5511999998888"}
	repo := &testutil.MockCobRepository{FindResult: expected}

	b := bo.NewGetCobBO(repo)
	result, err := b.Execute(bo.GetCobInput{TxID: "txid123"})

	require.NoError(t, err)
	assert.Equal(t, "txid123", result.Cob.TxID)
}

func TestGetCobBO_NotFound(t *testing.T) {
	repo := &testutil.MockCobRepository{FindResult: nil}

	b := bo.NewGetCobBO(repo)
	_, err := b.Execute(bo.GetCobInput{TxID: "naoexiste"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "não encontrada")
}

func TestGetCobBO_RepositoryError(t *testing.T) {
	repo := &testutil.MockCobRepository{FindError: testutil.ErrGeneric}

	b := bo.NewGetCobBO(repo)
	_, err := b.Execute(bo.GetCobInput{TxID: "txid123"})

	require.ErrorIs(t, err, testutil.ErrGeneric)
}
