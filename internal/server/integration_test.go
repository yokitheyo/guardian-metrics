package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	handlerpkg "github.com/yokitheyo/guardian-metrics/internal/server/handlers"
	"github.com/yokitheyo/guardian-metrics/internal/server/middleware"
	"github.com/yokitheyo/guardian-metrics/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestServerWithLogging(t *testing.T) {
	var buf bytes.Buffer

	config := zap.NewDevelopmentConfig()
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	))

	gin.SetMode(gin.TestMode)

	storage := storage.NewMemStorage()
	r := gin.New()

	r.Use(middleware.LoggingMiddleware(logger))

	r.POST("/update/:type/:name/:value", handlerpkg.UpdateMetricHandler(storage))
	r.GET("/value/:type/:name", handlerpkg.GetMetricValueHandler(storage))
	r.GET("/", handlerpkg.ListMetricsHandler(storage))

	req := httptest.NewRequest("POST", "/update/counter/testMetric/42", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest("GET", "/value/counter/testMetric", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	logOutput := buf.String()

	assert.Contains(t, logOutput, "POST")
	assert.Contains(t, logOutput, "/update/counter/testMetric/42")
	assert.Contains(t, logOutput, "GET")
	assert.Contains(t, logOutput, "/value/counter/testMetric")
	assert.Contains(t, logOutput, "HTTP request")
	assert.Contains(t, logOutput, "duration")
	assert.Contains(t, logOutput, "response_size")
	assert.Contains(t, logOutput, "status")

	logLines := bytes.Split(buf.Bytes(), []byte("\n"))
	nonEmptyLines := 0
	for _, line := range logLines {
		if len(bytes.TrimSpace(line)) > 0 {
			nonEmptyLines++
		}
	}
	assert.GreaterOrEqual(t, nonEmptyLines, 3, "Should have at least 3 log entries")
}
