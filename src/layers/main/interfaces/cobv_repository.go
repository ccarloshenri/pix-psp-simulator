package interfaces

import "pix-psp-simulator/src/layers/main/models"

type CobVRepository interface {
	Save(cobv models.CobV) error
	FindByTxID(txid string) (*models.CobV, error)
	Update(cobv models.CobV) error
}
