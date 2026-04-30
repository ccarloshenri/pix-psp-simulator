package memory

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

const numShards = 64

type cobShard struct {
	mu   sync.RWMutex
	data map[string]models.Cob
}

type CobRepository struct {
	shards [numShards]cobShard
}

func NewCobRepository() *CobRepository {
	r := &CobRepository{}
	for i := range r.shards {
		r.shards[i].data = make(map[string]models.Cob)
	}
	return r
}

func (r *CobRepository) shard(key string) *cobShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return &r.shards[h.Sum32()%numShards]
}

func (r *CobRepository) Save(cob models.Cob) error {
	s := r.shard(cob.TxID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[cob.TxID] = cob
	return nil
}

func (r *CobRepository) FindByTxID(txid string) (*models.Cob, error) {
	s := r.shard(txid)
	s.mu.RLock()
	defer s.mu.RUnlock()
	cob, ok := s.data[txid]
	if !ok {
		return nil, nil
	}
	return &cob, nil
}

func (r *CobRepository) Update(cob models.Cob) error {
	s := r.shard(cob.TxID)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[cob.TxID]; !ok {
		return fmt.Errorf("cobrança não encontrada")
	}
	s.data[cob.TxID] = cob
	return nil
}

func (r *CobRepository) Delete(txid string) error {
	s := r.shard(txid)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[txid]; !ok {
		return fmt.Errorf("cobrança não encontrada")
	}
	delete(s.data, txid)
	return nil
}

func (r *CobRepository) FindAll(filters interfaces.CobFilters) ([]models.Cob, error) {
	var inicio, fim time.Time
	if filters.Inicio != "" {
		inicio, _ = time.Parse(time.RFC3339, filters.Inicio)
	}
	if filters.Fim != "" {
		fim, _ = time.Parse(time.RFC3339, filters.Fim)
	}

	result := make([]models.Cob, 0)
	for i := range r.shards {
		s := &r.shards[i]
		s.mu.RLock()
		for _, cob := range s.data {
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
		s.mu.RUnlock()
	}
	return result, nil
}
