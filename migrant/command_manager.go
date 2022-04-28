package migrant

import (
	"fmt"
	"os"
)

func Run(migrant Migrant) {
	args := os.Args

	log := migrant.Logger()

	var err error
	if len(args) > 1 {
		switch args[1] {
		case "up":
			if len(args) != 3 {
				log.Error("Up command format must be: go run migrate up 00000000000_migation_name ")
				return
			}
			err = migrant.UpConcreteMigration(args[2])
			break
		case "down":
			if len(args) != 3 {
				log.Error("Down command format must be: go run migrate down 00000000000_migation_name ")
				return
			}
			err = migrant.DownConcreteMigration(args[2])
			break
		case "make":
			if len(args) != 3 {
				log.Error("Make command format must be: go run migrate.go make my_new_migration_name")
				return
			}
			err = migrant.MakeFileMigration(args[2])
			break
		default:
			err = fmt.Errorf("Unknown command parameters: %+v", args[1:])
		}
	} else {
		err = migrant.UpMigrations()
	}

	if err != nil {
		log.Error(err.Error())
	}
}
