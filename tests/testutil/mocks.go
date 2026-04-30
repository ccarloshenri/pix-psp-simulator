package testutil

import (
	"errors"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

// ErrGeneric is a sentinel error for tests that need a generic failure.
var ErrGeneric = errors.New("something went wrong")

// MockCobRepository tracks calls and allows injection of controlled responses.
type MockCobRepository struct {
	SaveCalled  bool
	SavedCob    models.Cob
	SaveError   error
	FindResult  *models.Cob
	FindError   error
	UpdateError error
	DeleteError error
}

func (m *MockCobRepository) Save(cob models.Cob) error {
	m.SaveCalled = true
	m.SavedCob = cob
	return m.SaveError
}

func (m *MockCobRepository) FindByTxID(_ string) (*models.Cob, error) {
	return m.FindResult, m.FindError
}

func (m *MockCobRepository) Update(cob models.Cob) error {
	return m.UpdateError
}

func (m *MockCobRepository) Delete(_ string) error {
	return m.DeleteError
}

// MockCobVRepository tracks calls and allows injection of controlled responses.
type MockCobVRepository struct {
	SaveCalled  bool
	SavedCobV   models.CobV
	SaveError   error
	FindResult  *models.CobV
	FindError   error
	UpdateError error
}

func (m *MockCobVRepository) Save(cobv models.CobV) error {
	m.SaveCalled = true
	m.SavedCobV = cobv
	return m.SaveError
}

func (m *MockCobVRepository) FindByTxID(_ string) (*models.CobV, error) {
	return m.FindResult, m.FindError
}

func (m *MockCobVRepository) Update(_ models.CobV) error {
	return m.UpdateError
}

// MockPixRepository tracks calls and allows injection of controlled responses.
type MockPixRepository struct {
	SaveCalled          bool
	SavedPix            models.Pix
	SaveError           error
	FindByE2EIDResult   *models.Pix
	FindByE2EIDError    error
	FindByTxIDResult    []models.Pix
	FindAllResult       []models.Pix
	FindAllError        error
	AddDevolucaoError   error
	UpdateDevolucaoError error
	FindDevolucaoResult *models.Devolucao
	FindDevolucaoError  error
}

func (m *MockPixRepository) Save(pix models.Pix) error {
	m.SaveCalled = true
	m.SavedPix = pix
	return m.SaveError
}

func (m *MockPixRepository) FindByE2EID(_ string) (*models.Pix, error) {
	return m.FindByE2EIDResult, m.FindByE2EIDError
}

func (m *MockPixRepository) FindByTxID(_ string) ([]models.Pix, error) {
	return m.FindByTxIDResult, nil
}

func (m *MockPixRepository) FindAll(_ interfaces.PixFilters) ([]models.Pix, error) {
	return m.FindAllResult, m.FindAllError
}

func (m *MockPixRepository) AddDevolucao(_ string, _ models.Devolucao) error {
	return m.AddDevolucaoError
}

func (m *MockPixRepository) UpdateDevolucao(_ string, _ models.Devolucao) error {
	return m.UpdateDevolucaoError
}

func (m *MockPixRepository) FindDevolucao(_, _ string) (*models.Devolucao, error) {
	return m.FindDevolucaoResult, m.FindDevolucaoError
}

// MockIDGenerator returns deterministic values for tests.
type MockIDGenerator struct {
	TxIDResult  string
	E2EIDResult string
	RtrIDResult string
}

func (m *MockIDGenerator) GenerateTxID() string {
	return m.TxIDResult
}

func (m *MockIDGenerator) GenerateE2EID(_ string) string {
	return m.E2EIDResult
}

func (m *MockIDGenerator) GenerateRtrID(_ string) string {
	return m.RtrIDResult
}
