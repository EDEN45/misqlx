package migrant

import (
	"github.com/jmoiron/sqlx"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var pool migrationsPool

type migrationsPool struct {
	migrations map[string]Migration
	sync.Mutex
}

func init() {
	pool = migrationsPool{migrations: make(map[string]Migration)}
}

type Migration interface {
	Up(db *sqlx.Tx, log Logger) error
	Down(db *sqlx.Tx, log Logger) error
}

// Each migration file call this method in its init method
func RegisterMigration(migration Migration) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic("Fail invoke caller")
		return
	}
	migrationName := strings.Replace(filepath.Base(file), ".go", "", -1)

	pool.Lock()
	defer pool.Unlock()
	_, ok = pool.migrations[migrationName]
	if ok {
		panic("Migration with name : " + migrationName + " already exist")
	}
	pool.migrations[migrationName] = migration
}
