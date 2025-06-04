package store

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Metric struct {
	ID    string
	MType MetricType
	Value *float64
	Delta *int64
}

type Storage interface {
	UpdateMetric(m Metric) error
	GetAll() []Metric
}
