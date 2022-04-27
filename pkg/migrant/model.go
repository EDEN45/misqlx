package migrant

import "time"

const DefaultMigrationTableName = "migrations"

type migrationModel struct {
	Id        uint32     `db:"id"`
	Name      string     `db:"name"`
	CreatedAt *time.Time `db:"created_at"`
}
