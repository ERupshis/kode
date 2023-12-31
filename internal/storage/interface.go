package storage

import (
	"database/sql"
	"fmt"

	"github.com/erupshis/kode.git/internal/config"
	"github.com/erupshis/kode.git/internal/logger"

	_ "github.com/lib/pq"
)

type Manager interface {
	AddText(user string, text string) error
	GetTexts(user string) ([]string, error)
	Close()
}

func CreateRamStorage() (Manager, error) {
	return &Storage{make(map[string][]string)}, nil
}

func CreatePostgresDB(cfg *config.Config, log logger.BaseLogger) (Manager, error) {
	dataSrcName := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword)
	log.Info("[storage::CreatePostgresDB] Start open postgres DB with settings: %s", dataSrcName)
	db, err := sql.Open("postgres", dataSrcName)
	if err != nil {
		log.Info("[storage::CreatePostgresDB] Failed to open postgres DB: %v", err)
	}

	storageDB := &postgresDB{database: db, log: &log}
	err = storageDB.createDataBase(cfg)
	if err != nil {
		log.Info("[storage::CreatePostgresDB] Failed to open '%s' Database: %v", cfg.DbName, err)
	}

	return storageDB, err
}
