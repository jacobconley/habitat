package peripherals

import (
	log "github.com/sirupsen/logrus"
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

	log.Debug("Building all peripherals...")
	
	if err, _ := BuildCSS(); err != nil { 
		log.Error("SASS build failed")
		return err 
	}

	
	if err := BuildWebpack(); err != nil {
		log.Error("Webpack build failed")
		return err
	}

	return nil 
}