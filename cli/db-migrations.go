package cli

import (
	"errors"
	"strconv"

	"github.com/gobuffalo/pop/v5"
	"github.com/jacobconley/habitat/habconf"
	"github.com/rs/zerolog/log"
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


		conf := habconf.GetConfig()


		conn, err := conf.NewConnection()
		if err != nil { 
			log.Err(err).Msg("Could not establish connection")
			return err
		}

		conn.Open()
		defer conn.Close()
		

		mig, err := pop.NewFileMigrator(conf.GetDirMigrations(), conn)
		if err != nil { 
			return err
		}

		mig.SchemaPath = conf.GetDirDB()

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