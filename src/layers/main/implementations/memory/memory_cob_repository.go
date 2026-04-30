package memory

import (
	"fmt"
	"sync"

	"pix-psp-simulator/src/layers/main/models"
)

type CobRepository struct {
	mu   sync.RWMutex
	data map[string]models.Cob
}

func NewCobRepository() *CobRepository {
	return &CobRepository{data: make(map[string]models.Cob)}
}

func (r *CobRepository) Save(cob models.Cob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[cob.TxID] = cob
	return nil
}

func (r *CobRepository) FindByTxID(txid string) (*models.Cob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cob, ok := r.data[txid]
	if !ok {
		return nil, nil
	}
	return &cob, nil
}

func (r *CobRepository) Update(cob models.Cob) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[cob.TxID]; !ok {
		return fmt.Errorf("cobrança não encontrada")
	}
	r.data[cob.TxID] = cob
	return nil
}

func (r *CobRepository) Delete(txid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[txid]; !ok {
		return fmt.Errorf("cobrança não encontrada")
	}
	delete(r.data, txid)
	return nil
}
