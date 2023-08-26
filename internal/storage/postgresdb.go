package storage

import (
	"database/sql"
	"strings"

	"github.com/erupshis/kode.git/internal/logger"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type UserTexts struct {
	Username string   `json:"username"`
	Texts    []string `json:"texts"`
}

type postgresDB struct {
	database *sql.DB
	log      *logger.BaseLogger
}

func (db *postgresDB) createDataBaseIfNeed(dbName string) error {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)"
	err := db.database.QueryRow(query, strings.ToLower(dbName)).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		if _, err = db.database.Exec("CREATE DATABASE IF NOT EXISTS " + strings.ToLower(dbName)); err != nil {
			return err
		}
	}

	if _, err = db.database.Exec(`CREATE TABLE IF NOT EXISTS user_texts (username VARCHAR PRIMARY KEY, texts TEXT[]);`); err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) AddText(username, text string) error {
	_, err := db.database.Exec("UPDATE user_texts SET texts = array_append(texts, $1) WHERE username = $2", text, username)
	return err
}

func (db *postgresDB) GetTexts(username string) ([]string, error) {
	var userTexts UserTexts
	row := db.database.QueryRow("SELECT username, texts FROM user_texts WHERE username = $1", username)
	err := row.Scan(&userTexts.Username, pq.Array(&userTexts.Texts))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return userTexts.Texts, nil
}

func (db *postgresDB) Close() {
	db.database.Close()
}
