package agent

type MetricsCollector interface {
	CollectMetrics() map[string]float64
}

type MetricsSender interface {
	SendMetrics(metrics map[string]float64) error
}
