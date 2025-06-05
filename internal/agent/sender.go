package agent

import (
	"fmt"
	"net/http"
	"strconv"
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
		metricType := "gauge"
		if name == "PollCount" {
			metricType = "counter"
		}

		url := fmt.Sprintf("%s/update/%s/%s/%s",
			s.serverAddress,
			metricType,
			name,
			strconv.FormatFloat(value, 'f', -1, 64),
		)

		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "text/plain")

		resp, err := s.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	return nil
}
