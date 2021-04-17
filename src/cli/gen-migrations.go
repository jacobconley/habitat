package cli

import (
	"context"
	"errors"
	habitat "habitat/src"
	"regexp"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/pop/v5/genny/fizz/cempty"
	"github.com/gobuffalo/pop/v5/genny/fizz/ctable"
	"github.com/spf13/cobra"
)

// h/t Buffalo/Soda/Fizz
// cf. https://github.com/gobuffalo/soda/blob/master/cmd/generate/fizz_cmd.go

var genMigrationCmd = &cobra.Command { 
	Use: 	"migration <name> [attrs]",
	Short: 	"Generates a new pair of up/down Fizz migrations with the given name",

	RunE: func(cmd * cobra.Command, args []string) error { 

		name := ""
		if len(args) > 0 { 
			name = args[0]
			if !(regexp.MustCompile(`^[\w_]+$`).MatchString(name)) { 
				return errors.New("For now, the migration name must consist only of alphanumeric characters and underscores - feel free to yell on the issues page if you want")
			}
		} else { 
			return errors.New("Required table name")
		}

		var ( 
			atts attrs.Attrs 
			err error
		)
		if len(args) > 1 { 
			atts, err = attrs.ParseArgs( args[1:]... )
			if err != nil { 
				return err
			}
		}


		runner := genny.WetRunner(context.Background())
		
		config, err := habitat.GetConfig()
		if err != nil { 
			return err 
		}
		path := config.GetDirMigrations()


		var g * genny.Generator
		if len(atts) == 0 { 

			g, err = cempty.New(&cempty.Options{
				Name: name,
				Path: path,
				Type: "fizz",
			})

		} else { 

			g, err = ctable.New(&ctable.Options{
				TableName: name,
				Path: path,
				Type: "fizz",
				Attrs: atts,
			})
			
		}
		
		if err != nil { 
			return err
		}
		runner.With(g)
		return runner.Run()
	},
}

func init(){ 
	cmdGen.AddCommand(genMigrationCmd)
}