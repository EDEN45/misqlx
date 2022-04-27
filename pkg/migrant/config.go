package migrant

type Config struct {
	Log           Logger
	MigrationsDir string
	TableName     string
}
