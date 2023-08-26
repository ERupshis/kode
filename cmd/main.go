package main

import (
	"net/http"

	"github.com/erupshis/kode.git/internal/config"
	"github.com/erupshis/kode.git/internal/controller"
	"github.com/erupshis/kode.git/internal/logger"
	"github.com/erupshis/kode.git/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Parse()

	log := logger.CreateZapLogger(cfg.LogLevel)
	defer log.Sync()

	storage := storage.CreateRamStorage()

	controller := controller.Create(log, storage)

	router := chi.NewRouter()
	router.Mount("/", controller.Route())

	log.Info("Server started with Host setting: %s", cfg.Host)
	if err := http.ListenAndServe(cfg.Host, router); err != nil {
		panic(err)
	}
}
