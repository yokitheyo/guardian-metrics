package store

import (
	"errors"
	"sync"
)

type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateMetric(m Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch m.MType {
	case Gauge:
		if m.Value != nil {
			s.gauges[m.ID] = *m.Value
		}
	case Counter:
		if m.Delta != nil {
			s.counters[m.ID] += *m.Delta
		}
	default:
		return errors.New("invalid metric type")
	}
	return nil
}

func (s *MemStorage) GetAll() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []Metric
	for id, val := range s.gauges {
		v := val
		result = append(result, Metric{ID: id, MType: Gauge, Value: &v})
	}
	for id, delta := range s.counters {
		d := delta
		result = append(result, Metric{ID: id, MType: Counter, Delta: &d})
	}
	return result
}

func (s *MemStorage) GetGauge(name string) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.gauges[name]
	return val, ok
}

func (s *MemStorage) GetCounter(name string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.counters[name]
	return val, ok
}
