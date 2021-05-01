package cli

import (
	"context"
	"errors"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/v2"
	"github.com/jacobconley/habitat/habconf"
	"github.com/spf13/cobra"

	"github.com/gobuffalo/pop/v5/genny/fizz/ctable"
	gmodel "github.com/gobuffalo/pop/v5/genny/model"
)

// h/t Buffalo/Pop
// cf. https://github.com/gobuffalo/soda/blob/master/cmd/generate/model_cmd.go

var optSkipMigration = false 

var genModelCmd = &cobra.Command { 
	Use: "model <name> [attrs]",
	Short: "Generates a Pop model with the given name and attributes",

	RunE: func(cmd * cobra.Command, args []string) error { 

		name := ""
		if len(args) > 0 { 
			name = args[0] 
		} else { 
			return errors.New("Model name is required")
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

		conf := habconf.GetConfig()
		runner := genny.WetRunner(context.Background())
		
		
		g, err := gmodel.New(  &gmodel.Options { 
			Name: 	name, 
			Attrs:  atts, 

			Encoding: 				"json",
			Path: 					conf.GetDirModels(),
			ForceDefaultID:  		true,
			ForceDefaultTimestamps:	true,
		})
		if err != nil { 
			return err
		}


		err = runner.With(g)
		if err != nil { 
			return err 
		}




		if !optSkipMigration { 

			// If we add support for SQL migrations [Issue #18], we will need a Translator here - see https://github.com/gobuffalo/soda/blob/master/cmd/generate/model_cmd.go#L82


			gm, err := ctable.New(&ctable.Options{
				TableName: 		name,
				Attrs: 			atts,

				Type: 						"fizz",
				Path:						conf.GetDirMigrations(),
				ForceDefaultTimestamps:  	true,
				ForceDefaultID:  			true,
			})
			if err != nil { 
				return err
			}

			err = runner.With(gm)
			if err != nil { 
				return err 
			}

		}


		return runner.Run()
	},
}


func init() { 
	genModelCmd.Flags().BoolVarP(&optSkipMigration, "skip-migration", "s", false, "Skip creating a migration for this model")
	cmdGen.AddCommand(genModelCmd)
}