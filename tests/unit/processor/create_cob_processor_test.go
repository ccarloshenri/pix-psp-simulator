package processor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pix-psp-simulator/src/layers/main/bo"
	"pix-psp-simulator/src/layers/main/models"
	"pix-psp-simulator/src/layers/main/processor"
	"pix-psp-simulator/tests/testutil"
)

func newCreateCobProcessor() *processor.CreateCobProcessor {
	repo := &testutil.MockCobRepository{}
	gen := &testutil.MockIDGenerator{TxIDResult: "txid-auto"}
	b := bo.NewCreateCobBO(repo, gen)
	return processor.NewCreateCobProcessor(b)
}

func TestCreateCobProcessor_MissingChave(t *testing.T) {
	p := newCreateCobProcessor()
	_, err := p.Process("", processor.CreateCobRequest{
		Valor: processor.CobValorRequest{Original: "100.00"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "chave")
}

func TestCreateCobProcessor_MissingValor(t *testing.T) {
	p := newCreateCobProcessor()
	_, err := p.Process("", processor.CreateCobRequest{
		Chave: "+5511999998888",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "valor")
}

func TestCreateCobProcessor_Success(t *testing.T) {
	p := newCreateCobProcessor()
	resp, err := p.Process("", processor.CreateCobRequest{
		Chave: "+5511999998888",
		Valor: processor.CobValorRequest{Original: "100.00"},
	})
	require.NoError(t, err)
	assert.Equal(t, "txid-auto", resp.Cob.TxID)
}

func TestCreateCobProcessor_WithTxID(t *testing.T) {
	p := newCreateCobProcessor()
	txid := "abcdefghij1234567890abcdef" // 26 chars alphanumeric
	resp, err := p.Process(txid, processor.CreateCobRequest{
		Chave: "+5511999998888",
		Valor: processor.CobValorRequest{Original: "200.00"},
	})
	require.NoError(t, err)
	assert.Equal(t, txid, resp.Cob.TxID)
}

func TestCreateCobProcessor_InvalidTxID(t *testing.T) {
	p := newCreateCobProcessor()
	_, err := p.Process("curto", processor.CreateCobRequest{
		Chave: "+5511999998888",
		Valor: processor.CobValorRequest{Original: "100.00"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "txid")
}

func TestCreateCobProcessor_InvalidValorFormat(t *testing.T) {
	p := newCreateCobProcessor()
	_, err := p.Process("", processor.CreateCobRequest{
		Chave: "+5511999998888",
		Valor: processor.CobValorRequest{Original: "100"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "valor.original")
}

func TestCreateCobProcessor_DevedorCPFandCNPJ(t *testing.T) {
	p := newCreateCobProcessor()
	_, err := p.Process("", processor.CreateCobRequest{
		Chave: "+5511999998888",
		Valor: processor.CobValorRequest{Original: "100.00"},
		Devedor: &models.Devedor{
			CPF:  "12345678909",
			CNPJ: "12345678000195",
			Nome: "Test",
		},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cpf")
}
