package memory

import (
	"fmt"
	"sync"
	"time"

	"pix-psp-simulator/src/layers/main/interfaces"
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

func (r *CobVRepository) FindAll(filters interfaces.CobVFilters) ([]models.CobV, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var inicio, fim time.Time
	if filters.Inicio != "" {
		inicio, _ = time.Parse(time.RFC3339, filters.Inicio)
	}
	if filters.Fim != "" {
		fim, _ = time.Parse(time.RFC3339, filters.Fim)
	}

	result := make([]models.CobV, 0)
	for _, cobv := range r.data {
		if filters.Status != "" && cobv.Status != filters.Status {
			continue
		}
		if filters.DataDeVencimento != "" && cobv.Calendario.DataDeVencimento != filters.DataDeVencimento {
			continue
		}
		if !inicio.IsZero() && cobv.Calendario.Criacao.Before(inicio) {
			continue
		}
		if !fim.IsZero() && cobv.Calendario.Criacao.After(fim) {
			continue
		}
		result = append(result, cobv)
	}
	return result, nil
}
