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

	//storageRam, _ := storage.CreateRamStorage()
	//serverController := controller.Create(log, storageRam)
	storageDB, err := storage.CreatePostgresDB(&cfg, log)
	if err != nil {
		panic(err)
	}
	defer storageDB.Close()
	serverController := controller.Create(log, storageDB)

	router := chi.NewRouter()
	router.Mount("/", serverController.Route())

	log.Info("[main] Server started with Host setting: %s", cfg.Host)
	if err := http.ListenAndServe(cfg.Host, router); err != nil {
		panic(err)
	}
}
