package migrant

import "fmt"

const (
	checkMigrationTableQuery  = "SELECT table_name FROM information_schema.tables WHERE table_name = '%s'"
	createMigrationTableQuery = `
		CREATE TABLE IF NOT EXISTS %s (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);`

	findMigrationsQuery  = "SELECT * FROM %s WHERE name = any($1)"
	findFirstQuery       = "SELECT * FROM %s WHERE name = $1"
	insertMigrationQuery = "INSERT INTO %s(name) values($1)"
	deleteMigrationQuery = "DELETE FROM %s WHERE name = $1"
)

func (m *migrant) checkTableQuery() string {
	return fmt.Sprintf(checkMigrationTableQuery, m.TableName)
}

func (m *migrant) createTableQuery() string {
	return fmt.Sprintf(createMigrationTableQuery, m.TableName)
}

func (m *migrant) findMigrationsQuery() string {
	return fmt.Sprintf(findMigrationsQuery, m.TableName)
}

func (m *migrant) findFirstQuery() string {
	return fmt.Sprintf(findFirstQuery, m.TableName)
}

func (m *migrant) insertMigrationQuery() string {
	return fmt.Sprintf(insertMigrationQuery, m.TableName)
}

func (m *migrant) deleteMigrationQuery() string {
	return fmt.Sprintf(deleteMigrationQuery, m.TableName)
}
