package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handlerpkg "github.com/yokitheyo/guardian-metrics/internal/server/handlers"
	storagepkg "github.com/yokitheyo/guardian-metrics/internal/storage"
)

func TestRunServer_Gin(t *testing.T) {
	storage := storagepkg.NewMemStorage()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := gin.New()
		r2.POST("/update/:type/:name/:value", handlerpkg.UpdateMetricHandler(storage))
		r2.ServeHTTP(w, r)
	}))
	defer ts.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(ts.URL+"/update/counter/testMetric/42", "text/plain", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	metrics := storage.GetAll()
	require.Len(t, metrics, 1)
	assert.Equal(t, "testMetric", metrics[0].ID)
	assert.Equal(t, storagepkg.Counter, metrics[0].MType)
	require.NotNil(t, metrics[0].Delta)
	assert.Equal(t, int64(42), *metrics[0].Delta)
}
