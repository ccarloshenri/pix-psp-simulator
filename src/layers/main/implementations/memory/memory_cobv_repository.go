package memory

import (
	"fmt"
	"sync"

	"pix-psp-simulator/src/layers/main/models"
)

type CobVRepository struct {
	mu   sync.RWMutex
	data map[string]models.CobV
}

func NewCobVRepository() *CobVRepository {
	return &CobVRepository{data: make(map[string]models.CobV)}
}

func (r *CobVRepository) Save(cobv models.CobV) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[cobv.TxID] = cobv
	return nil
}

func (r *CobVRepository) FindByTxID(txid string) (*models.CobV, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cobv, ok := r.data[txid]
	if !ok {
		return nil, nil
	}
	return &cobv, nil
}

func (r *CobVRepository) Update(cobv models.CobV) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[cobv.TxID]; !ok {
		return fmt.Errorf("cobrança com vencimento não encontrada")
	}
	r.data[cobv.TxID] = cobv
	return nil
}
