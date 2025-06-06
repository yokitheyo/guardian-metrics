package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	storagepkg "github.com/yokitheyo/guardian-metrics/internal/storage"
)

type MockStorage struct {
	metrics []storagepkg.Metric
}

func (m *MockStorage) UpdateMetric(metric storagepkg.Metric) error {
	m.metrics = append(m.metrics, metric)
	return nil
}

func (m *MockStorage) GetAll() []storagepkg.Metric {
	return m.metrics
}

func (m *MockStorage) GetGauge(name string) (float64, bool) {
	for _, metric := range m.metrics {
		if metric.ID == name && metric.MType == storagepkg.Gauge && metric.Value != nil {
			return *metric.Value, true
		}
	}
	return 0, false
}

func (m *MockStorage) GetCounter(name string) (int64, bool) {
	for _, metric := range m.metrics {
		if metric.ID == name && metric.MType == storagepkg.Counter && metric.Delta != nil {
			return *metric.Delta, true
		}
	}
	return 0, false
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

				var m storagepkg.Metric
				m.ID = name
				m.MType = storagepkg.MetricType(mType)

				switch m.MType {
				case storagepkg.Gauge:
					val, err := parseFloat(value)
					if err != nil {
						c.String(http.StatusBadRequest, "invalid gauge value")
						return
					}
					m.Value = &val
				case storagepkg.Counter:
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

func TestUpdateMetricHandler_Gin(t *testing.T) {
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
			r.POST("/update/:type/:name/:value", UpdateMetricHandler(storage))

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

func TestGetMetricValueHandler_Gin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storage := &MockStorage{}
	val := 42.5
	storage.metrics = append(storage.metrics, storagepkg.Metric{ID: "gaugeMetric", MType: storagepkg.Gauge, Value: &val})
	cnt := int64(10)
	storage.metrics = append(storage.metrics, storagepkg.Metric{ID: "counterMetric", MType: storagepkg.Counter, Delta: &cnt})

	r := gin.New()
	r.GET("/value/:type/:name", GetMetricValueHandler(storage))

	t.Run("existing gauge", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/value/gauge/gaugeMetric", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "42.5")
	})
	t.Run("existing counter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/value/counter/counterMetric", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "10")
	})
	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/value/gauge/unknown", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "metric not found")
	})
}

// func TestListMetricsHandler_Gin(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	storage := &MockStorage{}
// 	val := 1.23
// 	cnt := int64(7)
// 	storage.metrics = append(storage.metrics, store.Metric{ID: "gauge1", MType: store.Gauge, Value: &val})
// 	storage.metrics = append(storage.metrics, store.Metric{ID: "counter1", MType: store.Counter, Delta: &cnt})

// 	r := gin.New()
// 	r.GET("/", ListMetricsHandler(storage))

// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
// 	rr := httptest.NewRecorder()
// 	r.ServeHTTP(rr, req)

// 	// 	assert.Equal(t, http.StatusOK, rr.Code)
// 	assert.Contains(t, rr.Body.String(), "gauge1")
// 	assert.Contains(t, rr.Body.String(), "counter1")
// 	assert.Contains(t, rr.Body.String(), "1.230000")
// 	assert.Contains(t, rr.Body.String(), ">7<")
// }
