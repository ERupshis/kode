package storage

import (
	"database/sql"

	"github.com/erupshis/kode.git/internal/logger"
	_ "github.com/lib/pq"
)

type UserTexts struct {
	Username string   `json:"username"`
	Texts    []string `json:"texts"`
}

type postgresDB struct {
	db  *sql.DB
	log *logger.BaseLogger
}

// TODO: create new database, create new table
// TODO: open database and table
// TODO: handle with it.
func (db *postgresDB) AddText(username, text string) error {

	return nil
}

func (db *postgresDB) GetTexts(username string) ([]string, error) {

	return []string{}, nil
}
