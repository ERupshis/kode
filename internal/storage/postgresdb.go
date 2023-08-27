package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/erupshis/kode.git/internal/config"
	"github.com/erupshis/kode.git/internal/logger"
	"github.com/lib/pq"
)

const schemaName = "notes"

type postgresDB struct {
	database *sql.DB
	log      *logger.BaseLogger
}

func (db *postgresDB) createDataBase(cfg *config.Config) error {
	if err := db.createDataBaseIfNeed(cfg); err != nil {
		return err
	}

	if err := db.openDB(cfg); err != nil {
		return err
	}

	if err := db.createSchemaIfNeed(cfg); err != nil {
		return err
	}

	if err := db.createTableIfNeed(); err != nil {
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
	if _, err = db.database.Exec(grantDatabaseSQL); err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) openDB(cfg *config.Config) error {
	if db.database != nil {
		db.Close()
	}

	dataSrcName := fmt.Sprintf(" user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbUser, cfg.DbPassword, strings.ToLower(cfg.DbName))

	var err error
	db.database, err = sql.Open("postgres", dataSrcName)
	return err
}

func (db *postgresDB) createSchemaIfNeed(cfg *config.Config) error {
	createSchemaSQL := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schemaName)
	if _, err := db.database.Exec(createSchemaSQL); err != nil {
		return err
	}

	if _, err := db.database.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)); err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) createTableIfNeed() error {
	createTableSql := `CREATE TABLE IF NOT EXISTS notes 
		(time TIMESTAMP PRIMARY KEY, username TEXT, note TEXT);`

	if _, err := db.database.Exec(createTableSql); err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) AddText(username, text string) error {
	addTextSQL := fmt.Sprintf(`INSERT INTO %s.notes (time, username, note) 
		VALUES (clock_timestamp(), '%s', '%s');`,
		schemaName, username, text)

	if _, err := db.database.Exec(addTextSQL); err != nil {
		return err
	}

	return nil
}

func (db *postgresDB) GetTexts(username string) ([]string, error) {
	getTextsSQL := fmt.Sprintf(`SELECT ARRAY(
    	SELECT note FROM %s.notes WHERE username = '%s' ORDER BY time ASC);`,
		schemaName, username)

	row := db.database.QueryRow(getTextsSQL)
	var textArray pq.StringArray
	err := row.Scan(&textArray)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return textArray, nil
}

func (db *postgresDB) Close() {
	db.database.Close()
}
