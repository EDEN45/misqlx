package migrations

import (
	"github.com/EDEN45/misqlx/pkg/migrant"
	"github.com/jmoiron/sqlx"
)

func init() {
	migrant.RegisterMigration(&migrationFirst{})
}

type migrationFirst struct{}

func (m *migrationFirst) Up(tx *sqlx.Tx, log migrant.Logger) error {
	return nil
}

func (m *migrationFirst) Down(tx *sqlx.Tx, log migrant.Logger) error {
	return nil
}
