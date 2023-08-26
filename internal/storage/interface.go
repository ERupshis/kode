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
}

func CreateRamStorage() (Manager, error) {
	return &Storage{make(map[string][]string)}, nil
}

func CreatePostgresDB(cfg *config.Config, log logger.BaseLogger) (Manager, error) {
	dataSrcName := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable", cfg.DbUser, cfg.DbPassword)
	db, err := sql.Open("postgres", dataSrcName)
	if err != nil {
		log.Info("[storage::CreatePostgresDB] failed to open postgres DB with user: %s, password: %s", cfg.DbUser, cfg.DbPassword)
	}
	return &postgresDB{db: db, log: &log}, err
}
