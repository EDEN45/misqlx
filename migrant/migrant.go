package migrant

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"strings"
	"time"
)

const DefaultMigrationDir = "migrations"

type Migrant interface {
	Logger() Logger
	UpMigrations() error
	UpConcreteMigration(name string) error
	DownConcreteMigration(name string) error
	MakeFileMigration(name string) error
}

type migrant struct {
	*Config
	db *sqlx.DB
}

func NewMigrant(db *sqlx.DB, cfg *Config) Migrant {
	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.Log == nil {
		cfg.Log = NewLogger()
	}
	if cfg.MigrationsDir == "" {
		cfg.MigrationsDir = DefaultMigrationDir
	}
	if cfg.TableName == "" {
		cfg.TableName = DefaultMigrationTableName
	}
	if db == nil {
		cfg.Log.Error("Empty dp connection")
		return nil
	}

	return &migrant{
		Config: cfg,
		db:     db,
	}
}

func (m *migrant) Logger() Logger {
	return m.Log
}

// UpMigrations - Up all migrations
func (m *migrant) UpMigrations() error {
	m.Log.Info("Start migrations")

	m.checkMigrationTable()

	newMigrations := m.getNewMigrations()

	successCnt := 0
	for _, migration := range newMigrations {
		if migration.Id == 0 {
			tx, _ := m.db.Beginx()
			if err := pool.migrations[migration.Name].Up(tx, m.Log); err != nil {
				tx.Rollback()
				return fmt.Errorf("up migration: %+v, err: %+v", migration.Name, err)
			}

			if _, err := tx.Exec(m.insertMigrationQuery(), migration.Name); err != nil {
				tx.Rollback()
				return fmt.Errorf("save migration: %v, err: %+v", migration.Name, err)
			}
			tx.Commit()
			m.Log.Info("success: %+v", migration.Name)
			successCnt++
		}
	}

	if successCnt > 0 {
		m.Log.Info("All migrations are done success!")
	} else {
		m.Log.Info("Nothing to migrate.")
	}

	return nil
}

func (m *migrant) UpConcreteMigration(name string) error {
	mig, ok := pool.migrations[name]
	if !ok {
		return errors.New("Does not exist migration with name: " + name)
	}

	tx, err := m.db.Beginx()
	if err := mig.Up(tx, m.Log); err != nil {
		tx.Rollback()
		return err
	}

	migrationModel := migrationModel{}
	err = tx.Get(&migrationModel, m.findFirstQuery(), name)
	if err != sql.ErrNoRows && err != nil {
		return err
	}

	if migrationModel.Id == 0 {
		if _, err := tx.Exec(m.insertMigrationQuery(), name); err != nil {
			tx.Rollback()
			return fmt.Errorf("save migration: %v, err: %+v", name, err)
		}
	}
	tx.Commit()

	return nil
}

func (m *migrant) DownConcreteMigration(name string) error {
	mig, ok := pool.migrations[name]
	if !ok {
		return errors.New("Does not exist migration with name: " + name)
	}

	m.Log.Info("Start down migration: %s", name)

	tx, err := m.db.Beginx()
	if err := mig.Down(tx, m.Log); err != nil {
		tx.Rollback()
		return err
	}

	migrationModel := migrationModel{}
	err = tx.Get(&migrationModel, m.findFirstQuery(), name)
	if err != sql.ErrNoRows && err != nil {
		return err
	}

	if migrationModel.Id != 0 {
		if _, err := tx.Exec(m.deleteMigrationQuery(), name); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	m.Log.Info("Migration: %s is down finish!", name)

	return nil
}

func (m *migrant) MakeFileMigration(name string) error {
	migrationsPath := m.Config.MigrationsDir

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		m.Log.Info("Create new directory : %v", migrationsPath)
		if err := os.MkdirAll(migrationsPath, os.ModePerm); err != nil {
			return err
		}
	}

	err := checkFileExists(migrationsPath, name+".go")
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	realName := fmt.Sprintf("%d_%s.go", now, name)

	migrationPath := migrationsPath + "/" + realName

	f, err := os.Create(migrationPath)
	if err != nil {
		return fmt.Errorf("create migration file: %v", err)
	}

	partsName := strings.Split(name, "_")
	structName := "migration"
	for _, p := range partsName {
		structName += strings.Title(p)
	}

	partsDir := strings.Split(m.Config.MigrationsDir, "/")
	packageName := partsDir[len(partsDir)-1]

	tmpl, err := getTemplate()
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, map[string]interface{}{"struct_name": structName, "package": packageName, "file_name": realName})

	if err != nil {
		return err
	}

	m.Log.Info("migration file created: %v", realName)

	return nil
}
