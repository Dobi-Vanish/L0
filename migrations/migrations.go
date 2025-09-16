package migrations

import (
	"L0/pkg/errormsg"
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var EmbedMigrations embed.FS

// Apply applies all available migrations via Goose.
func Apply(db *sql.DB) error {
	goose.SetBaseFS(EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return errormsg.ErrSetDialect
	}

	if err := goose.Up(db, "."); err != nil {
		return errormsg.ErrApplyMigrations
	}

	return nil
}
