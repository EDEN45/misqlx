package migrations

import (
	migrant2 "github.com/EDEN45/misqlx/migrant"
	"github.com/jmoiron/sqlx"
)

func init() {
	migrant2.RegisterMigration(&migrationFirst{})
}

type migrationFirst struct{}

func (m *migrationFirst) Up(tx *sqlx.Tx, log migrant2.Logger) error {
	return nil
}

func (m *migrationFirst) Down(tx *sqlx.Tx, log migrant2.Logger) error {
	return nil
}
