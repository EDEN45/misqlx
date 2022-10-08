package main

import (
	"github.com/EDEN45/misqlx/example/config"
	migrant2 "github.com/EDEN45/misqlx/migrant"
	"github.com/jmoiron/sqlx"
	"path"
	"runtime"

	_ "github.com/EDEN45/misqlx/example/migrations"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.GetConfig()
	db, err := sqlx.Connect(cfg.Driver, cfg.GetDSN())
	if err != nil {
		panic(err)
	}

	_, file, _, _ := runtime.Caller(0)
	curDir := path.Dir(file)
	mg := migrant2.NewMigrant(db, &migrant2.Config{
		MigrationsDir: curDir + "/migrations",
		TableName:     "migrations",
		Log:           migrant2.NewLogger(),
	})

	migrant2.Run(mg)
}
