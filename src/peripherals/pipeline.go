package peripherals

import (
	habitat "habitat/src"
	"habitat/src/peripherals/loaders"
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


// GetLoadersFromConfig gets all peripheral loaders
func GetLoadersFromConfig(config * habitat.Config) []Loader { 
	return []Loader { 

		loaders.NewCSSFromConfig(config),

	}
}