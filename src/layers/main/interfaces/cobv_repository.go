package interfaces

import "pix-psp-simulator/src/layers/main/models"

type CobVFilters struct {
	Status           string
	DataDeVencimento string // YYYY-MM-DD
	Inicio           string // RFC3339
	Fim              string // RFC3339
}

type CobVRepository interface {
	Save(cobv models.CobV) error
	FindByTxID(txid string) (*models.CobV, error)
	FindAll(filters CobVFilters) ([]models.CobV, error)
	Update(cobv models.CobV) error
}
