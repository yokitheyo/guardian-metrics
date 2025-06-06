package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yokitheyo/guardian-metrics/internal/store"
)

func UpdateMetricHandler(storage store.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		mType := c.Param("type")
		name := c.Param("name")
		value := c.Param("value")

		var m store.Metric
		m.ID = name
		m.MType = store.MetricType(mType)

		switch m.MType {
		case store.Gauge:
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				c.String(http.StatusBadRequest, "invalid gauge value")
				return
			}
			m.Value = &val
		case store.Counter:
			delta, err := strconv.ParseInt(value, 10, 64)
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
	}
}

func GetMetricValueHandler(storage store.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		mType := c.Param("type")
		name := c.Param("name")
		var (
			found  bool
			result string
		)
		switch store.MetricType(mType) {
		case store.Gauge:
			var val float64
			val, found = storage.GetGauge(name)
			if found {
				result = fmt.Sprintf("%.2f", val)
			}
		case store.Counter:
			var val int64
			val, found = storage.GetCounter(name)
			if found {
				result = fmt.Sprintf("%d", val)
			}
		}
		if !found {
			c.String(http.StatusNotFound, "metric not found")
			return
		}
		c.String(http.StatusOK, result)
	}
}

func ListMetricsHandler(storage store.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics := storage.GetAll()
		tmpl := `<html><body><h1>Metrics</h1><table border="1"><tr><th>Name</th><th>Type</th><th>Value</th></tr>{{range .}}<tr><td>{{.ID}}</td><td>{{.MType}}</td><td>{{if eq .MType "gauge"}}{{with .Value}}{{printf "%f" .}}{{end}}{{else}}{{.Delta}}{{end}}</td></tr>{{end}}</table></body></html>`
		t, err := template.New("metrics").Parse(tmpl)
		if err != nil {
			c.String(http.StatusInternalServerError, "template error")
			return
		}
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		t.Execute(c.Writer, metrics)
	}
}
