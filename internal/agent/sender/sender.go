package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPSender struct {
	client        *http.Client
	serverAddress string
}

func NewHTTPSender(serverAddress string) *HTTPSender {
	return &HTTPSender{
		client:        &http.Client{},
		serverAddress: serverAddress,
	}
}

func (s *HTTPSender) SendMetrics(metrics map[string]float64) error {
	for name, value := range metrics {
		// Determine metric type
		var metricType string
		if name == "PollCount" {
			metricType = "counter"
		} else {
			metricType = "gauge"
		}

		// Create metric struct
		var metric map[string]interface{}
		if metricType == "counter" {
			// Convert float64 to int64 for counter
			delta := int64(value)
			metric = map[string]interface{}{
				"id":    name,
				"type":  metricType,
				"delta": delta,
			}
		} else {
			metric = map[string]interface{}{
				"id":    name,
				"type":  metricType,
				"value": value,
			}
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(metric)
		if err != nil {
			return fmt.Errorf("failed to marshal metric: %w", err)
		}

		// Create request
		url := fmt.Sprintf("%s/update/", s.serverAddress)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		// Send request
		resp, err := s.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}

		if resp.Body != nil {
			resp.Body.Close()
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	return nil
}
