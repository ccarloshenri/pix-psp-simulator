package memory

import (
	"fmt"
	"sync"
	"time"

	"pix-psp-simulator/src/layers/main/interfaces"
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

func (r *CobRepository) FindAll(filters interfaces.CobFilters) ([]models.Cob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var inicio, fim time.Time
	if filters.Inicio != "" {
		inicio, _ = time.Parse(time.RFC3339, filters.Inicio)
	}
	if filters.Fim != "" {
		fim, _ = time.Parse(time.RFC3339, filters.Fim)
	}

	result := make([]models.Cob, 0)
	for _, cob := range r.data {
		if filters.Status != "" && cob.Status != filters.Status {
			continue
		}
		if !inicio.IsZero() && cob.Calendario.Criacao.Before(inicio) {
			continue
		}
		if !fim.IsZero() && cob.Calendario.Criacao.After(fim) {
			continue
		}
		result = append(result, cob)
	}
	return result, nil
}
