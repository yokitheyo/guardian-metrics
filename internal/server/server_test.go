package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yokitheyo/guardian-metrics/internal/server/handler"
	"github.com/yokitheyo/guardian-metrics/internal/store"
)

func TestRunServer(t *testing.T) {
	storage := store.NewMemStorage()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handler.NewUpdateHandler(storage)
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodPost, server.URL+"/update/counter/testMetric/42", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	metrics := storage.GetAll()
	require.Len(t, metrics, 1)
	assert.Equal(t, "testMetric", metrics[0].ID)
	assert.Equal(t, store.Counter, metrics[0].MType)
	require.NotNil(t, metrics[0].Delta)
	assert.Equal(t, int64(42), *metrics[0].Delta)
}
