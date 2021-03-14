package e2e

import (
	habitat "habitat/src"
	"habitat/src/peripherals/loaders"
	"os"
	"testing"
)

//TODO: Initialization and cleanup

func TestCSS(t *testing.T) {
	os.Chdir("../../test-fixtures/userland")
	
	config, _ := habitat.GetConfig()
	loader := loaders.NewCSSFromConfig(config)

	if err := loader.Build(); err != nil { 
		t.Error("Build failed;", err) 
	}
	
}