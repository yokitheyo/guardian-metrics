package agent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockCollector мок для тестирования
type MockCollector struct {
	metrics map[string]float64
}

func (m *MockCollector) CollectMetrics() map[string]float64 {
	return m.metrics
}

// MockSender мок для тестирования
type MockSender struct {
	sentMetrics []map[string]float64
}

func (m *MockSender) SendMetrics(metrics map[string]float64) error {
	m.sentMetrics = append(m.sentMetrics, metrics)
	return nil
}

func TestAgent(t *testing.T) {
	// Создаем моки
	collector := &MockCollector{
		metrics: map[string]float64{
			"TestMetric": 42.0,
		},
	}
	sender := &MockSender{}

	// Создаем агента с короткими интервалами для тестирования
	a := NewAgent(
		collector,
		sender,
		100*time.Millisecond,
		200*time.Millisecond,
		"http://localhost:8080",
	)

	// Запускаем агента в горутине
	go a.Run()

	// Ждем некоторое время для сбора и отправки метрик
	time.Sleep(300 * time.Millisecond)

	// Проверяем, что метрики были отправлены
	require.NotEmpty(t, sender.sentMetrics, "No metrics were sent")

	// Проверяем содержимое отправленных метрик
	lastMetrics := sender.sentMetrics[len(sender.sentMetrics)-1]
	assert.Contains(t, lastMetrics, "TestMetric", "TestMetric should be present in metrics")
	assert.Contains(t, lastMetrics, "PollCount", "PollCount should be present in metrics")
	assert.Contains(t, lastMetrics, "RandomValue", "RandomValue should be present in metrics")

	// Проверяем значения метрик
	assert.Equal(t, 42.0, lastMetrics["TestMetric"], "TestMetric value should be 42.0")
	assert.Greater(t, lastMetrics["PollCount"], float64(0), "PollCount should be greater than 0")
	assert.Greater(t, lastMetrics["RandomValue"], float64(0), "RandomValue should be greater than 0")
}
