package main

import (
	"log"

	"github.com/yokitheyo/guardian-metrics/internal/server"
	"github.com/yokitheyo/guardian-metrics/internal/store"
)

func main() {
	storage := store.NewMemStorage()

	if err := server.RunServer(storage); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
