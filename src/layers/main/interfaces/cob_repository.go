package interfaces

import "pix-psp-simulator/src/layers/main/models"

type CobFilters struct {
	Status string
	Inicio string // RFC3339
	Fim    string // RFC3339
}

type CobRepository interface {
	Save(cob models.Cob) error
	FindByTxID(txid string) (*models.Cob, error)
	FindAll(filters CobFilters) ([]models.Cob, error)
	Update(cob models.Cob) error
	Delete(txid string) error
}
