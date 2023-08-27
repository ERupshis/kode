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

const schemaName = "notes"

type UserTexts struct {
	Username string   `json:"username"`
	Texts    []string `json:"texts"`
}

type postgresDB struct {
	database *sql.DB
	log      *logger.BaseLogger
}

func (db *postgresDB) createDataBaseIfNeedAndOpen(cfg *config.Config) error {
	if err := db.createDataBaseIfNeed(cfg); err != nil {
		return err
	}

	if err := db.createSchemaIfNeed(cfg); err != nil {
		return err
	}

	if err := db.createTableIfNeed(cfg); err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) createDataBaseIfNeed(cfg *config.Config) error {
	dbName := strings.ToLower(cfg.DbName)
	var exists bool
	checkDbSQL := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);"
	err := db.database.QueryRow(checkDbSQL, dbName).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		createDbSQL := fmt.Sprintf("CREATE DATABASE %s OWNER =  %s;", dbName, cfg.DbUser)
		if _, err = db.database.Exec(createDbSQL); err != nil {
			return err
		}
	}

	grantDatabaseSQL := fmt.Sprintf("GRANT ALL ON DATABASE %s to %s;", cfg.DbName, cfg.DbUser)
	if _, err := db.database.Exec(grantDatabaseSQL); err != nil {
		return err
	}

	db.database.Close()

	dataSrcName := fmt.Sprintf(" user=%s password=%s dbname=%s sslmode=disable", cfg.DbUser, cfg.DbPassword, dbName)
	db.database, err = sql.Open("postgres", dataSrcName)
	if err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) createSchemaIfNeed(cfg *config.Config) error {
	createSchemaSQL := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schemaName)
	if _, err := db.database.Exec(createSchemaSQL); err != nil {
		return err
	}

	grantSchemaSQL := fmt.Sprintf("GRANT ALL PRIVILEGES ON SCHEMA %s TO %s;", schemaName, cfg.DbUser)
	if _, err := db.database.Exec(grantSchemaSQL); err != nil {
		return err
	}

	if _, err := db.database.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)); err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) createTableIfNeed(cfg *config.Config) error {
	createTableSql := "CREATE TABLE IF NOT EXISTS user_texts (username VARCHAR PRIMARY KEY, texts TEXT[]);"
	if _, err := db.database.Exec(createTableSql); err != nil {
		return err
	}

	grantTableSQL := fmt.Sprintf("GRANT ALL PRIVILEGES ON TABLE %s TO %s;", "user_texts", cfg.DbUser)
	if _, err := db.database.Exec(grantTableSQL); err != nil {
		return err
	}

	grantSequencesSQL := fmt.Sprintf("GRANT SELECT, UPDATE, USAGE ON ALL SEQUENCES IN SCHEMA %s TO %s;", schemaName, cfg.DbUser)
	if _, err := db.database.Exec(grantSequencesSQL); err != nil {
		return err
	}

	grantFuncsSQL := fmt.Sprintf("GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA %s TO %s;", schemaName, cfg.DbUser)
	if _, err := db.database.Exec(grantFuncsSQL); err != nil {
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
