package server

import (
	"log"

	"github.com/gin-gonic/gin"
	handlerpkg "github.com/yokitheyo/guardian-metrics/internal/server/handlers"
	"github.com/yokitheyo/guardian-metrics/internal/server/middleware"
	"github.com/yokitheyo/guardian-metrics/internal/storage"
	"go.uber.org/zap"
)

func RunServer(storage storage.Storage, addr string, logger *zap.Logger) error {
	r := gin.Default()

	r.Use(middleware.LoggingMiddleware(logger))

	r.POST("/update/:type/:name/:value", handlerpkg.UpdateMetricHandler(storage))
	r.GET("/value/:type/:name", handlerpkg.GetMetricValueHandler(storage))
	r.GET("/", handlerpkg.ListMetricsHandler(storage))

	r.POST("/update/", handlerpkg.UpdateMetricJSONHandler(storage))
	r.POST("/value/", handlerpkg.GetMetricJSONHandler(storage))

	log.Println("starting server on", addr)
	return r.Run(addr)
}
