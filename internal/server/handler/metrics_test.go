package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yokitheyo/guardian-metrics/internal/store"
)

type MockStorage struct {
	metrics []store.Metric
}

func (m *MockStorage) UpdateMetric(metric store.Metric) error {
	m.metrics = append(m.metrics, metric)
	return nil
}

func (m *MockStorage) GetAll() []store.Metric {
	return m.metrics
}

func (m *MockStorage) GetGauge(name string) (float64, bool) {
	for _, metric := range m.metrics {
		if metric.ID == name && metric.MType == store.Gauge && metric.Value != nil {
			return *metric.Value, true
		}
	}
	return 0, false
}

func (m *MockStorage) GetCounter(name string) (int64, bool) {
	for _, metric := range m.metrics {
		if metric.ID == name && metric.MType == store.Counter && metric.Delta != nil {
			return *metric.Delta, true
		}
	}
	return 0, false
}

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid counter metric",
			method:         http.MethodPost,
			path:           "/update/counter/someMetric/527",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "valid gauge metric",
			method:         http.MethodPost,
			path:           "/update/gauge/someMetric/42.5",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			path:           "/update/counter/someMetric/527",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "missing metric name",
			method:         http.MethodPost,
			path:           "/update/counter//527",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "metric name not provided\n",
		},
		{
			name:           "missing metric value",
			method:         http.MethodPost,
			path:           "/update/counter/someMetric",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "metric value not provided\n",
		},
		{
			name:           "invalid metric type",
			method:         http.MethodPost,
			path:           "/update/invalid/someMetric/527",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid metric type\n",
		},
		{
			name:           "invalid counter value",
			method:         http.MethodPost,
			path:           "/update/counter/someMetric/invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid counter value\n",
		},
		{
			name:           "invalid gauge value",
			method:         http.MethodPost,
			path:           "/update/gauge/someMetric/invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid gauge value\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MockStorage{}

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Content-Type", "text/plain")

			rr := httptest.NewRecorder()

			handler := NewUpdateHandler(storage)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			assert.Equal(t, tt.expectedBody, rr.Body.String())

			if tt.expectedStatus == http.StatusOK {
				require.Len(t, storage.metrics, 1)
				metric := storage.metrics[0]
				assert.Equal(t, "someMetric", metric.ID)
			}
		})
	}
}

func TestUpdateHandler_Gin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid counter metric",
			method:         http.MethodPost,
			path:           "/update/counter/someMetric/527",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "valid gauge metric",
			method:         http.MethodPost,
			path:           "/update/gauge/someMetric/42.5",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "invalid metric type",
			method:         http.MethodPost,
			path:           "/update/invalid/someMetric/527",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid metric type",
		},
		{
			name:           "invalid counter value",
			method:         http.MethodPost,
			path:           "/update/counter/someMetric/invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid counter value",
		},
		{
			name:           "invalid gauge value",
			method:         http.MethodPost,
			path:           "/update/gauge/someMetric/invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid gauge value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MockStorage{}
			r := gin.New()
			r.POST("/update/:type/:name/:value", func(c *gin.Context) {
				mType := c.Param("type")
				name := c.Param("name")
				value := c.Param("value")

				var m store.Metric
				m.ID = name
				m.MType = store.MetricType(mType)

				switch m.MType {
				case store.Gauge:
					val, err := parseFloat(value)
					if err != nil {
						c.String(http.StatusBadRequest, "invalid gauge value")
						return
					}
					m.Value = &val
				case store.Counter:
					delta, err := parseInt(value)
					if err != nil {
						c.String(http.StatusBadRequest, "invalid counter value")
						return
					}
					m.Delta = &delta
				default:
					c.String(http.StatusBadRequest, "invalid metric type")
					return
				}

				if err := storage.UpdateMetric(m); err != nil {
					c.String(http.StatusInternalServerError, "failed to update")
					return
				}
				c.String(http.StatusOK, "OK")
			})

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.Header.Set("Content-Type", "text/plain")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())

			if tt.expectedStatus == http.StatusOK {
				require.Len(t, storage.metrics, 1)
				metric := storage.metrics[0]
				assert.Equal(t, "someMetric", metric.ID)
			}
		})
	}
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
