package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yokitheyo/guardian-metrics/internal/store"
)

func RunServer(addr string, storage store.Storage) error {
	r := gin.Default()

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

	r.GET("/value/:type/:name", func(c *gin.Context) {
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
				result = fmt.Sprintf("%f", val)
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
	})

	r.GET("/", func(c *gin.Context) {
		metrics := storage.GetAll()
		tmpl := `<html><body><h1>Metrics</h1><table border="1"><tr><th>Name</th><th>Type</th><th>Value</th></tr>{{range .}}<tr><td>{{.ID}}</td><td>{{.MType}}</td><td>{{if eq .MType "gauge"}}{{printf "%f" .Value}}{{else}}{{.Delta}}{{end}}</td></tr>{{end}}</table></body></html>`
		t, err := template.New("metrics").Parse(tmpl)
		if err != nil {
			c.String(http.StatusInternalServerError, "template error")
			return
		}
		c.Status(http.StatusOK)
		c.Header("Content-Type", "text/html; charset=utf-8")
		t.Execute(c.Writer, metrics)
	})

	log.Println("starting server on", addr)
	return r.Run(addr)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
