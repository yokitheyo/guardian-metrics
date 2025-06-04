package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/yokitheyo/guardian-metrics/internal/store"
)

func NewUpdateHandler(storage store.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")
		if len(parts) != 3 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		mType, name, rawVal := parts[0], parts[1], parts[2]
		if name == "" {
			http.Error(w, "missing metric name", http.StatusNotFound)
			return
		}

		var m store.Metric
		m.ID = name
		m.MType = store.MetricType(mType)

		switch m.MType {
		case store.Gauge:
			val, err := strconv.ParseFloat(rawVal, 64)
			if err != nil {
				http.Error(w, "invalid gauge value", http.StatusBadRequest)
				return
			}
			m.Value = &val
		case store.Counter:
			delta, err := strconv.ParseInt(rawVal, 10, 64)
			if err != nil {
				http.Error(w, "invalid counter value", http.StatusBadRequest)
				return
			}
			m.Delta = &delta

		default:
			http.Error(w, "invalid metric type", http.StatusBadRequest)
		}
		if err := storage.UpdateMetric(m); err != nil {
			http.Error(w, "failed to update", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
