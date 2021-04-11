package cli

import (
	"errors"
	habitat "habitat/src"

	"github.com/gobuffalo/pop/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Hat tip to the Buffalo/Pop/Soda folks, as always
// cf. https://github.com/gobuffalo/soda/blob/master/cmd/migrate.go

var migrateCmd = &cobra.Command{ 
	Use: 	"migrate",
	Short: 	"Runs all pending migrations on your database",

	RunE: func(cmd *cobra.Command, args []string) error { 
		if len(args) > 0 { 
			return errors.New("This command does not accept any arguments")
		}

		config, err := habitat.GetConfig()
		if err != nil { 
			return err 
		}

		conn, err := config.NewConnection()
		if err != nil { 
			log.Error("Could not establish connection: ", err)
			return err
		}

		mig, err := pop.NewFileMigrator(config.GetDirMigrations(), conn)
		if err != nil { 
			return err
		}

		mig.SchemaPath = config.GetDirDB()

		return mig.Up()
	},
}


func init() { 
	cmdDB.AddCommand(migrateCmd)
}