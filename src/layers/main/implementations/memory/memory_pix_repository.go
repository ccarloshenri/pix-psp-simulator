package memory

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"pix-psp-simulator/src/layers/main/interfaces"
	"pix-psp-simulator/src/layers/main/models"
)

type pixShard struct {
	mu   sync.RWMutex
	data map[string]models.Pix
}

type PixRepository struct {
	shards [numShards]pixShard
}

func NewPixRepository() *PixRepository {
	r := &PixRepository{}
	for i := range r.shards {
		r.shards[i].data = make(map[string]models.Pix)
	}
	return r
}

func (r *PixRepository) shard(key string) *pixShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return &r.shards[h.Sum32()%numShards]
}

func (r *PixRepository) Save(pix models.Pix) error {
	s := r.shard(pix.EndToEndID)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[pix.EndToEndID] = pix
	return nil
}

func (r *PixRepository) FindByE2EID(e2eid string) (*models.Pix, error) {
	s := r.shard(e2eid)
	s.mu.RLock()
	defer s.mu.RUnlock()
	pix, ok := s.data[e2eid]
	if !ok {
		return nil, nil
	}
	return &pix, nil
}

func (r *PixRepository) FindByTxID(txid string) ([]models.Pix, error) {
	var result []models.Pix
	for i := range r.shards {
		s := &r.shards[i]
		s.mu.RLock()
		for _, pix := range s.data {
			if pix.TxID == txid {
				result = append(result, pix)
			}
		}
		s.mu.RUnlock()
	}
	return result, nil
}

func (r *PixRepository) FindAll(filters interfaces.PixFilters) ([]models.Pix, error) {
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
	for i := range r.shards {
		s := &r.shards[i]
		s.mu.RLock()
		for _, pix := range s.data {
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
		s.mu.RUnlock()
	}
	return result, nil
}

func (r *PixRepository) AddDevolucao(e2eid string, dev models.Devolucao) error {
	s := r.shard(e2eid)
	s.mu.Lock()
	defer s.mu.Unlock()
	pix, ok := s.data[e2eid]
	if !ok {
		return fmt.Errorf("pagamento não encontrado")
	}
	pix.Devolucoes = append(pix.Devolucoes, dev)
	s.data[e2eid] = pix
	return nil
}

func (r *PixRepository) UpdateDevolucao(e2eid string, dev models.Devolucao) error {
	s := r.shard(e2eid)
	s.mu.Lock()
	defer s.mu.Unlock()
	pix, ok := s.data[e2eid]
	if !ok {
		return fmt.Errorf("pagamento não encontrado")
	}
	for i, d := range pix.Devolucoes {
		if d.ID == dev.ID {
			pix.Devolucoes[i] = dev
			s.data[e2eid] = pix
			return nil
		}
	}
	return fmt.Errorf("devolução não encontrada")
}

func (r *PixRepository) FindDevolucao(e2eid, devID string) (*models.Devolucao, error) {
	s := r.shard(e2eid)
	s.mu.RLock()
	defer s.mu.RUnlock()
	pix, ok := s.data[e2eid]
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
