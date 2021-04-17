package cli

import (
	"errors"
	habitat "habitat/src"
	"strconv"

	"github.com/gobuffalo/pop/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Hat tip to the Buffalo/Pop/Soda folks, as always
// cf. https://github.com/gobuffalo/soda/blob/master/cmd/migrate.go

var migrateCmd = &cobra.Command{ 
	Use: 	"migrate ([up] | down <number of steps>)",
	Short: 	"Runs all pending migrations on your database",

	RunE: func(cmd *cobra.Command, args []string) error { 
		
		var up = true
		var steps = 0 
		if len(args) > 0 { 

			if args[0] == "down" { 
				up = false

				if len(args) == 2 { 
					var err error
					steps, err = strconv.Atoi(args[1])
					if err != nil { 
						return err
					}
				} else { 
					return errors.New("Expected number of steps")
				}

			} else if !(len(args) == 1 && args[0] == "up") { 
				return errors.New("Invalid argument")
			}

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

		if up { 
			return mig.Up()
		} else { 
			return mig.Down(steps)
		}
	},
}


func init() { 
	cmdDB.AddCommand(migrateCmd)
}