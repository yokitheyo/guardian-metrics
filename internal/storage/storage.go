package storage

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta *int64     `json:"delta,omitempty"`
	Value *float64   `json:"value,omitempty"`
}

type Storage interface {
	UpdateMetric(m Metric) error
	GetAll() []Metric
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
}
