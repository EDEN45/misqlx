package migrant

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"path"
	"runtime"
	"sort"
	"strings"
)

// check or create table to register successful migrations
func (m *migrant) checkMigrationTable() {
	tableName := ""
	err := m.db.Get(&tableName, m.checkTableQuery())
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if tableName != m.TableName {
		m.Log.Info("Init table: %v", m.Config.TableName)
		_, err := m.db.Exec(fmt.Sprintf(createMigrationTableQuery, DefaultMigrationTableName))
		if err != nil {
			panic(err)
		}
	}
}

// Finds not yet completed migration files
func (m *migrant) getNewMigrations() []migrationModel {

	var names []string
	for k, _ := range pool.migrations {
		names = append(names, k)
	}

	sort.Strings(names)

	step := 20 // limit
	result := make([]migrationModel, 0)
	existMigrations := make(map[string]bool)
	for i := 0; i < len(names); {

		i += step
		var chunkNames []string
		if i <= len(names) {
			chunkNames = names[i-step : i]
		} else {
			chunkNames = names[i-step:]
		}

		var rows []*migrationModel
		err := m.db.Select(&rows, m.findMigrationsQuery(), pq.Array(chunkNames))
		if err != nil {
			panic(err)
		}

		for _, row := range rows {
			existMigrations[row.Name] = true
		}
	}

	for _, name := range names {
		if _, ok := existMigrations[name]; !ok {
			model := migrationModel{}
			model.Name = name
			result = append(result, model)
		}
	}

	return result
}

// сheck the existence of a file in the directory with migrations
func checkFileExists(dir string, name string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		split := strings.Split(f.Name(), "_")

		if name == strings.Join(split[1:], "_") {
			return fmt.Errorf("File %v already exists in dir: %v", name, dir)
		}
	}

	return nil
}

func getTemplate() (*template.Template, error) {

	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return nil, fmt.Errorf("Template caller")
	}

	tmpl, err := template.ParseFiles(path.Dir(filename) + "/../assets/template")
	if err != nil {
		return nil, fmt.Errorf("parse template : %v", err)
	}

	return tmpl, nil
}
