package memory

import (
	"fmt"
	"sync"
	"time"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type PixRepository struct {
	mu   sync.RWMutex
	data map[string]models.Pix
}

func NewPixRepository() *PixRepository {
	return &PixRepository{data: make(map[string]models.Pix)}
}

func (r *PixRepository) Save(pix models.Pix) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[pix.EndToEndID] = pix
	return nil
}

func (r *PixRepository) FindByE2EID(e2eid string) (*models.Pix, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pix, ok := r.data[e2eid]
	if !ok {
		return nil, nil
	}
	return &pix, nil
}

func (r *PixRepository) FindByTxID(txid string) ([]models.Pix, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []models.Pix
	for _, pix := range r.data {
		if pix.TxID == txid {
			result = append(result, pix)
		}
	}
	return result, nil
}

func (r *PixRepository) FindAll(filters interfaces.PixFilters) ([]models.Pix, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var inicio, fim time.Time
	if filters.Inicio != "" {
		parsed, err := time.Parse(time.RFC3339, filters.Inicio)
		if err != nil {
			return nil, fmt.Errorf("formato de inicio inválido, esperado RFC3339: %w", err)
		}
		inicio = parsed
	}
	if filters.Fim != "" {
		parsed, err := time.Parse(time.RFC3339, filters.Fim)
		if err != nil {
			return nil, fmt.Errorf("formato de fim inválido, esperado RFC3339: %w", err)
		}
		fim = parsed
	}

	var result []models.Pix
	for _, pix := range r.data {
		if filters.TxID != "" && pix.TxID != filters.TxID {
			continue
		}
		if !inicio.IsZero() && pix.Horario.Before(inicio) {
			continue
		}
		if !fim.IsZero() && pix.Horario.After(fim) {
			continue
		}
		result = append(result, pix)
	}
	return result, nil
}

func (r *PixRepository) AddDevolucao(e2eid string, dev models.Devolucao) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	pix, ok := r.data[e2eid]
	if !ok {
		return fmt.Errorf("pagamento não encontrado")
	}
	pix.Devolucoes = append(pix.Devolucoes, dev)
	r.data[e2eid] = pix
	return nil
}

func (r *PixRepository) UpdateDevolucao(e2eid string, dev models.Devolucao) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	pix, ok := r.data[e2eid]
	if !ok {
		return fmt.Errorf("pagamento não encontrado")
	}
	for i, d := range pix.Devolucoes {
		if d.ID == dev.ID {
			pix.Devolucoes[i] = dev
			r.data[e2eid] = pix
			return nil
		}
	}
	return fmt.Errorf("devolução não encontrada")
}

func (r *PixRepository) FindDevolucao(e2eid, devID string) (*models.Devolucao, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pix, ok := r.data[e2eid]
	if !ok {
		return nil, nil
	}
	for _, d := range pix.Devolucoes {
		if d.ID == devID {
			return &d, nil
		}
	}
	return nil, nil
}
