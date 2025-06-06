package agent

import (
	"log"
	"time"
)

type Agent struct {
	collector      MetricsCollector
	sender         MetricsSender
	pollInterval   time.Duration
	reportInterval time.Duration
	serverAddress  string
}

func NewAgent(collector MetricsCollector, sender MetricsSender, pollInterval, reportInterval time.Duration, serverAddress string) *Agent {
	return &Agent{
		collector:      collector,
		sender:         sender,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		serverAddress:  serverAddress,
	}
}

func (a *Agent) Run() {
	pollTicker := time.NewTicker(a.pollInterval)
	reportTicker := time.NewTicker(a.reportInterval)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	metrics := make(map[string]float64)
	pollCount := int64(0)

	for {
		select {
		case <-pollTicker.C:
			runtimeMetrics := a.collector.CollectMetrics()
			for k, v := range runtimeMetrics {
				metrics[k] = v
			}
			pollCount++
			metrics["PollCount"] = float64(pollCount)
			metrics["RandomValue"] = float64(time.Now().UnixNano())

		case <-reportTicker.C:
			if err := a.sender.SendMetrics(metrics); err != nil {
				log.Printf("failed to send metrics: %v", err)
			}
		}
	}
}
