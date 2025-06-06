package main

import (
	"log"

	"github.com/yokitheyo/guardian-metrics/internal/config"
	"github.com/yokitheyo/guardian-metrics/internal/server"
	"github.com/yokitheyo/guardian-metrics/internal/storage"
)

func main() {
	cfg := config.LoadServerConfig()
	storage := storage.NewMemStorage()

	if err := server.RunServer(storage, cfg.Address); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
