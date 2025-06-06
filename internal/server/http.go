package server

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yokitheyo/guardian-metrics/internal/server/handler"
	"github.com/yokitheyo/guardian-metrics/internal/store"
)

func RunServer(storage store.Storage) error {
	addr := flag.String("a", "localhost:8080", "address for HTTP server")
	flag.Parse()

	r := gin.Default()

	r.POST("/update/:type/:name/:value", handler.UpdateMetricHandler(storage))
	r.GET("/value/:type/:name", handler.GetMetricValueHandler(storage))
	r.GET("/", handler.ListMetricsHandler(storage))
	log.Println("starting server on", *addr)
	return r.Run(*addr)
}
