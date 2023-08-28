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

	_, err = db.Exec(fmt.Sprintf(`DROP DATABASE %s;`, cfg.DbName))
	require.NoError(t, err)

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
