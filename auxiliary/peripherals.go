package auxiliary

import (
	"github.com/rs/zerolog/log"
)

//Pipeline manages the processing of project files as defined by the config
// type Pipeline struct {

// 	Files 	[]string
// 	RecursiveDirs 	[]string

// }

// var pipeline *Pipeline
// func GetPipeline() *Pipeline {

// 	if pipeline != nil {
// 		return pipeline
// 	}

// 	log.Debug("Initializing pipeline")
// }

// Loader a type of thing that can be processed - CSS, JSS, Etc
// CSS composes SASS rather than any sort of inheritance shenanigans
type Loader interface {


	Build() error
	Watch() error 
	
}




func BuildAll() error { 

	log.Debug().Msg("Building all auxiliary...")
	
	if err, _ := BuildCSS(); err != nil { 
		log.Error().Msg("SASS build failed")
		return err 
	}

	
	if err := BuildWebpack(); err != nil {
		log.Error().Msg("Webpack build failed")
		return err
	}

	return nil 
}