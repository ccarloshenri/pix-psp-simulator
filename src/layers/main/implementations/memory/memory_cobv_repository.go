package memory

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type cobvShard struct {
	mu   sync.RWMutex
	data map[string]models.CobV
}

type CobVRepository struct {
	shards [numShards]cobvShard
}

func NewCobVRepository() *CobVRepository {
	r := &CobVRepository{}
	for i := range r.shards {
		r.shards[i].data = make(map[string]models.CobV)
	}
	return r
}

func (r *CobVRepository) shard(key string) *cobvShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return &r.shards[h.Sum32()%numShards]
}

func (r *CobVRepository) Save(cobv models.CobV) error {
	s := r.shard(cobv.TxID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[cobv.TxID] = cobv
	return nil
}

func (r *CobVRepository) FindByTxID(txid string) (*models.CobV, error) {
	s := r.shard(txid)
	s.mu.RLock()
	defer s.mu.RUnlock()
	cobv, ok := s.data[txid]
	if !ok {
		return nil, nil
	}
	return &cobv, nil
}

func (r *CobVRepository) Update(cobv models.CobV) error {
	s := r.shard(cobv.TxID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[cobv.TxID]; !ok {
		return fmt.Errorf("cobrança com vencimento não encontrada")
	}
	s.data[cobv.TxID] = cobv
	return nil
}

func (r *CobVRepository) FindAll(filters interfaces.CobVFilters) ([]models.CobV, error) {
	var inicio, fim time.Time
	if filters.Inicio != "" {
		inicio, _ = time.Parse(time.RFC3339, filters.Inicio)
	}
	if filters.Fim != "" {
		fim, _ = time.Parse(time.RFC3339, filters.Fim)
	}

	result := make([]models.CobV, 0)
	for i := range r.shards {
		s := &r.shards[i]
		s.mu.RLock()
		for _, cobv := range s.data {
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
		s.mu.RUnlock()
	}
	return result, nil
}
