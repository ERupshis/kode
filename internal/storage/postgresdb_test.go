//go:build ignoreSQLtest
// +build ignoreSQLtest

package storage

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/erupshis/kode.git/internal/config"
	"github.com/erupshis/kode.git/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_postgresDB_AddTextGetTexts(t *testing.T) {
	cfg := config.Config{
		DbHost:     "localhost",
		DbName:     "test_db",
		DbUser:     "postgres",
		DbPassword: "postgres",
	}

	dataSrcName := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword)
	db, err := sql.Open("postgres", dataSrcName)
	require.NoError(t, err)

	var exists bool
	checkDbSQL := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);"
	err = db.QueryRow(checkDbSQL, cfg.DbName).Scan(&exists)
	require.NoError(t, err)
	if exists {
		_, err = db.Exec(fmt.Sprintf(`DROP DATABASE %s;`, cfg.DbName))
		require.NoError(t, err)
	}

	log := logger.CreateZapLogger(cfg.LogLevel)
	storageDB := postgresDB{db, &log}
	storageDB.createDataBase(&cfg)

	type args struct {
		username string
		text     string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "add correct text",
			args: args{username: "asd", text: "correct"},
			want: []string{"correct"},
		},
		{
			name: "add 2nd correct text",
			args: args{username: "asd", text: "correct another text"},
			want: []string{"correct", "correct another text"},
		},
		{
			name: "add incorrect text",
			args: args{username: "asd", text: "incorect"},
			want: []string{"correct", "correct another text", "incorect"},
		},
		{
			name: "add correct text by another user",
			args: args{username: "zxc", text: "looks correct"},
			want: []string{"looks correct"},
		},
		{
			name: "add empty text by another user",
			args: args{username: "zxc", text: ""},
			want: []string{"looks correct", ""},
		},
		{
			name: "add incorrect text by another user",
			args: args{username: "zxc", text: "incorect"},
			want: []string{"looks correct", "", "incorect"},
		},
		{
			name: "add first empty text by another user",
			args: args{username: "o", text: ""},
			want: []string{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storageDB.AddText(tt.args.username, tt.args.text)
			require.NoError(t, err)

			notes, err := storageDB.GetTexts(tt.args.username)
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, notes)
		})
	}
}

func Test_postgresDB_createDataBaseIfNeed(t *testing.T) {
	cfg := config.Config{
		DbHost:     "localhost",
		DbName:     "test_db1",
		DbUser:     "postgres",
		DbPassword: "postgres",
	}

	dataSrcName := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword)
	db, err := sql.Open("postgres", dataSrcName)
	require.NoError(t, err)

	var exists bool
	checkDbSQL := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);"
	err = db.QueryRow(checkDbSQL, cfg.DbName).Scan(&exists)
	require.NoError(t, err)
	if exists {
		_, err = db.Exec(fmt.Sprintf(`DROP DATABASE %s;`, cfg.DbName))
		require.NoError(t, err)
	}

	log := logger.CreateZapLogger(cfg.LogLevel)
	storageDB := postgresDB{db, &log}

	//CREATE MISSING DB.
	err = storageDB.createDataBaseIfNeed(&cfg)
	require.NoError(t, err)

	err = storageDB.database.QueryRow(checkDbSQL, cfg.DbName).Scan(&exists)
	require.NoError(t, err)

	//CREATE ALREADY EXISTING DB.
	err = storageDB.createDataBaseIfNeed(&cfg)
	require.NoError(t, err)

	err = storageDB.database.QueryRow(checkDbSQL, cfg.DbName).Scan(&exists)
	require.NoError(t, err)
}

func Test_postgresDB_openDB(t *testing.T) {
	cfg := config.Config{
		DbHost:     "localhost",
		DbName:     "test_db2",
		DbUser:     "postgres",
		DbPassword: "postgres",
	}

	dataSrcName := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword)
	db, err := sql.Open("postgres", dataSrcName)
	require.NoError(t, err)

	var exists bool
	checkDbSQL := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);"
	err = db.QueryRow(checkDbSQL, cfg.DbName).Scan(&exists)
	require.NoError(t, err)
	if exists {
		_, err = db.Exec(fmt.Sprintf(`DROP DATABASE %s;`, cfg.DbName))
		require.NoError(t, err)
	}

	log := logger.CreateZapLogger(cfg.LogLevel)
	storageDB := postgresDB{db, &log}

	//OPEN MISSING DB
	err = storageDB.openDB(&cfg)
	assert.Error(t, err)

	//CREATE MISSING DB.
	err = storageDB.createDataBaseIfNeed(&cfg)
	require.NoError(t, err)

	err = storageDB.openDB(&cfg)
	require.NoError(t, err)

}

func Test_postgresDB_isExistDB(t *testing.T) {
	cfg := config.Config{
		DbHost:     "localhost",
		DbName:     "test_db3",
		DbUser:     "postgres",
		DbPassword: "postgres",
	}

	dataSrcName := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword)
	db, err := sql.Open("postgres", dataSrcName)
	require.NoError(t, err)

	var exists bool
	checkDbSQL := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);"
	err = db.QueryRow(checkDbSQL, cfg.DbName).Scan(&exists)
	require.NoError(t, err)
	if exists {
		_, err = db.Exec(fmt.Sprintf(`DROP DATABASE %s;`, cfg.DbName))
		require.NoError(t, err)
	}

	log := logger.CreateZapLogger(cfg.LogLevel)
	storageDB := postgresDB{db, &log}

	//CHECK MISSING DB
	exists, err = storageDB.isExistDB(&cfg)
	assert.NoError(t, err)
	assert.True(t, !exists)

	//CHECK EXISTING DB.
	err = storageDB.createDataBaseIfNeed(&cfg)
	require.NoError(t, err)

	exists, err = storageDB.isExistDB(&cfg)
	assert.NoError(t, err)
	assert.True(t, exists)

}

func Test_postgresDB_createSchemaIfNeed(t *testing.T) {
	cfg := config.Config{
		DbHost:     "localhost",
		DbName:     "test_db4",
		DbUser:     "postgres",
		DbPassword: "postgres",
	}

	dataSrcName := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword)
	db, err := sql.Open("postgres", dataSrcName)
	require.NoError(t, err)

	var exists bool
	checkDbSQL := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);"
	err = db.QueryRow(checkDbSQL, cfg.DbName).Scan(&exists)
	require.NoError(t, err)
	if exists {
		_, err = db.Exec(fmt.Sprintf(`DROP DATABASE %s;`, cfg.DbName))
		require.NoError(t, err)
	}

	log := logger.CreateZapLogger(cfg.LogLevel)
	storageDB := postgresDB{db, &log}

	err = storageDB.createDataBaseIfNeed(&cfg)
	require.NoError(t, err)

	err = storageDB.openDB(&cfg)
	require.NoError(t, err)

	//CHECK MISSING SCHEMA
	var schemaExists bool
	checkSchemaSQL := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = '%s');`, schemaName)
	err = storageDB.database.QueryRow(checkSchemaSQL).Scan(&schemaExists)
	require.NoError(t, err)
	assert.True(t, !schemaExists)

	//CHECK EXISTING SCHEMA
	err = storageDB.createSchemaIfNeed()
	assert.NoError(t, err)

	err = storageDB.database.QueryRow(checkSchemaSQL).Scan(&schemaExists)
	require.NoError(t, err)
	assert.True(t, schemaExists)
}

func Test_postgresDB_createTableIfNeed(t *testing.T) {
	cfg := config.Config{
		DbHost:     "localhost",
		DbName:     "test_db5",
		DbUser:     "postgres",
		DbPassword: "postgres",
	}

	dataSrcName := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword)
	db, err := sql.Open("postgres", dataSrcName)
	require.NoError(t, err)

	var exists bool
	checkDbSQL := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1);"
	err = db.QueryRow(checkDbSQL, cfg.DbName).Scan(&exists)
	require.NoError(t, err)
	if exists {
		_, err = db.Exec(fmt.Sprintf(`DROP DATABASE %s;`, cfg.DbName))
		require.NoError(t, err)
	}

	log := logger.CreateZapLogger(cfg.LogLevel)
	storageDB := postgresDB{db, &log}

	err = storageDB.createDataBaseIfNeed(&cfg)
	require.NoError(t, err)

	err = storageDB.openDB(&cfg)
	require.NoError(t, err)

	err = storageDB.createSchemaIfNeed()
	assert.NoError(t, err)

	//CHECK MISSING TABLE
	var tableExists bool
	checkSchemaSQL := fmt.Sprintf(`SELECT EXISTS (SELECT FROM pg_tables 
                      						WHERE  schemaname = '%s' 
                      						AND tablename  = '%s' );`,
		schemaName, tableName)

	err = storageDB.database.QueryRow(checkSchemaSQL).Scan(&tableExists)
	require.NoError(t, err)
	assert.True(t, !tableExists)

	//CHECK EXISTING SCHEMA
	err = storageDB.createTableIfNeed()
	assert.NoError(t, err)

	err = storageDB.database.QueryRow(checkSchemaSQL).Scan(&tableExists)
	require.NoError(t, err)
	assert.True(t, tableExists)
}
