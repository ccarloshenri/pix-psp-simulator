package interfaces

import "pix-psp-simulator/src/layers/main/models"

type PixFilters struct {
	TxID  string
	Inicio string // RFC3339
	Fim   string // RFC3339
}

type PixRepository interface {
	Save(pix models.Pix) error
	FindByE2EID(e2eid string) (*models.Pix, error)
	FindByTxID(txid string) ([]models.Pix, error)
	FindAll(filters PixFilters) ([]models.Pix, error)
	AddDevolucao(e2eid string, dev models.Devolucao) error
	UpdateDevolucao(e2eid string, dev models.Devolucao) error
	FindDevolucao(e2eid, devID string) (*models.Devolucao, error)
}
