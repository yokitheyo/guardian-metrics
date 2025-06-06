package server

import (
	"log"

	"github.com/gin-gonic/gin"
	handlerpkg "github.com/yokitheyo/guardian-metrics/internal/server/handlers"
	"github.com/yokitheyo/guardian-metrics/internal/storage"
)

func RunServer(storage storage.Storage, addr string) error {
	r := gin.Default()

	r.POST("/update/:type/:name/:value", handlerpkg.UpdateMetricHandler(storage))
	r.GET("/value/:type/:name", handlerpkg.GetMetricValueHandler(storage))
	r.GET("/", handlerpkg.ListMetricsHandler(storage))
	log.Println("starting server on", addr)
	return r.Run(addr)
}
