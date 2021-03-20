package main

import (
	"os"
)

func main() {

	os.Chdir("test-fixtures/userland")
	
	// config, _ := habitat.GetConfig()
	// loader := loaders.NewCSSFromConfig(config)

	// if err := loader.Build(); err != nil { 
	// 	log.Error("Build failed;", err) 
	// }
}