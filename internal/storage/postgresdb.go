package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/erupshis/kode.git/internal/config"
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

func (db *postgresDB) createDataBaseIfNeedAndOpen(cfg *config.Config) error {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)"
	err := db.database.QueryRow(query, strings.ToLower(cfg.DbName)).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		if _, err = db.database.Exec("CREATE DATABASE IF NOT EXISTS " + strings.ToLower(cfg.DbName)); err != nil {
			return err
		}
	}

	db.database.Close()
	dataSrcName := fmt.Sprintf(" user=%s password=%s dbname=%s sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbName)
	requiredDB, err := sql.Open("postgres", dataSrcName)
	if err != nil {
		return err
	}

	db.database = requiredDB
	_, err = db.database.Exec(`
		CREATE TABLE IF NOT EXISTS user_texts (
			username VARCHAR PRIMARY KEY,
			texts TEXT[]
		);
	`)

	return nil
}

func (db *postgresDB) AddText(username, text string) error {
	_, err := db.database.Exec("UPDATE user_texts SET texts = array_append(texts, $1) WHERE username = $2", text, username)
	return err
}

func (db *postgresDB) GetTexts(username string) ([]string, error) {
	var userTexts UserTexts
	row := db.database.QueryRow("SELECT username, texts FROM user_texts WHERE username = $1", username)
	var textArray pq.StringArray // Use pq.StringArray to scan array of strings
	err := row.Scan(&userTexts.Username, &textArray)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	userTexts.Texts = textArray
	return userTexts.Texts, nil
}

func (db *postgresDB) Close() {
	db.database.Close()
}
