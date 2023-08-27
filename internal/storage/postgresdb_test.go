package storage

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/erupshis/kode.git/internal/logger"
	"github.com/stretchr/testify/assert"
)

func Test_postgresDB_AddText(t *testing.T) {
	type fields struct {
		database *sql.DB
		log      *logger.BaseLogger
	}
	type args struct {
		username string
		text     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &postgresDB{
				database: tt.fields.database,
				log:      tt.fields.log,
			}
			tt.wantErr(t, db.AddText(tt.args.username, tt.args.text), fmt.Sprintf("AddText(%v, %v)", tt.args.username, tt.args.text))
		})
	}
}

func Test_postgresDB_GetTexts(t *testing.T) {
	type fields struct {
		database *sql.DB
		log      *logger.BaseLogger
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &postgresDB{
				database: tt.fields.database,
				log:      tt.fields.log,
			}
			got, err := db.GetTexts(tt.args.username)
			if !tt.wantErr(t, err, fmt.Sprintf("GetTexts(%v)", tt.args.username)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetTexts(%v)", tt.args.username)
		})
	}
}
