package server

import (
	"log"
	"net/http"

	"github.com/yokitheyo/guardian-metrics/internal/server/handler"
	"github.com/yokitheyo/guardian-metrics/internal/store"
)

func RunServer(addr string, storage store.Storage) error {
	mux := http.NewServeMux()
	mux.Handle("/update/", handler.NewUpdateHandler(storage))

	log.Println("starting server on", addr)
	return http.ListenAndServe(addr, mux)
}
