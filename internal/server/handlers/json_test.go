package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	storagepkg "github.com/yokitheyo/guardian-metrics/internal/storage"
)

func TestUpdateMetricJSONHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		metric         storagepkg.Metric
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid gauge metric",
			metric: storagepkg.Metric{
				ID:    "testGauge",
				MType: storagepkg.Gauge,
				Value: func() *float64 { v := 42.5; return &v }(),
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "valid counter metric",
			metric: storagepkg.Metric{
				ID:    "testCounter",
				MType: storagepkg.Counter,
				Delta: func() *int64 { v := int64(10); return &v }(),
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "invalid metric type",
			metric: storagepkg.Metric{
				ID:    "testMetric",
				MType: "invalid",
				Value: func() *float64 { v := 42.5; return &v }(),
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "gauge without value",
			metric: storagepkg.Metric{
				ID:    "testGauge",
				MType: storagepkg.Gauge,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "counter without delta",
			metric: storagepkg.Metric{
				ID:    "testCounter",
				MType: storagepkg.Counter,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MockStorage{}
			r := gin.New()
			r.POST("/update/", UpdateMetricJSONHandler(storage))

			jsonData, err := json.Marshal(tt.metric)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.expectError {
				var responseMetric storagepkg.Metric
				err := json.Unmarshal(rr.Body.Bytes(), &responseMetric)
				require.NoError(t, err)
				assert.Equal(t, tt.metric.ID, responseMetric.ID)
				assert.Equal(t, tt.metric.MType, responseMetric.MType)
			}
		})
	}
}

func TestGetMetricJSONHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	storage := &MockStorage{}
	gaugeVal := 42.5
	counterVal := int64(10)
	storage.metrics = []storagepkg.Metric{
		{ID: "testGauge", MType: storagepkg.Gauge, Value: &gaugeVal},
		{ID: "testCounter", MType: storagepkg.Counter, Delta: &counterVal},
	}

	r := gin.New()
	r.POST("/value/", GetMetricJSONHandler(storage))

	tests := []struct {
		name           string
		requestMetric  storagepkg.Metric
		expectedStatus int
		expectError    bool
		expectedValue  interface{}
	}{
		{
			name: "get existing gauge",
			requestMetric: storagepkg.Metric{
				ID:    "testGauge",
				MType: storagepkg.Gauge,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedValue:  42.5,
		},
		{
			name: "get existing counter",
			requestMetric: storagepkg.Metric{
				ID:    "testCounter",
				MType: storagepkg.Counter,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedValue:  int64(10),
		},
		{
			name: "get non-existing metric",
			requestMetric: storagepkg.Metric{
				ID:    "nonExisting",
				MType: storagepkg.Gauge,
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name: "invalid metric type",
			requestMetric: storagepkg.Metric{
				ID:    "testGauge",
				MType: "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "missing ID",
			requestMetric: storagepkg.Metric{
				MType: storagepkg.Gauge,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.requestMetric)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.expectError {
				var responseMetric storagepkg.Metric
				err := json.Unmarshal(rr.Body.Bytes(), &responseMetric)
				require.NoError(t, err)
				assert.Equal(t, tt.requestMetric.ID, responseMetric.ID)
				assert.Equal(t, tt.requestMetric.MType, responseMetric.MType)

				if tt.requestMetric.MType == storagepkg.Gauge {
					require.NotNil(t, responseMetric.Value)
					assert.Equal(t, tt.expectedValue, *responseMetric.Value)
				} else if tt.requestMetric.MType == storagepkg.Counter {
					require.NotNil(t, responseMetric.Delta)
					assert.Equal(t, tt.expectedValue, *responseMetric.Delta)
				}
			}
		})
	}
}
