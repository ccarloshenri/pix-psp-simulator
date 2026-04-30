package interfaces

import "pix-psp-simulator/src/layers/main/models"

type CobRepository interface {
	Save(cob models.Cob) error
	FindByTxID(txid string) (*models.Cob, error)
	Update(cob models.Cob) error
	Delete(txid string) error
}
