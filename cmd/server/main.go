package main

import (
	"log"

	"github.com/yokitheyo/guardian-metrics/internal/config"
	"github.com/yokitheyo/guardian-metrics/internal/server"
	"github.com/yokitheyo/guardian-metrics/internal/storage"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadServerConfig()
	storage := storage.NewMemStorage()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	if err := server.RunServer(storage, cfg.Address, logger); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
