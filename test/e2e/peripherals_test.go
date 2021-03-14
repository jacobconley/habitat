package e2e

import (
	habitat "habitat/src"
	"habitat/src/peripherals/loaders"
	"io/ioutil"
	"os"
	"testing"
)

func TestCSS(t *testing.T) {
	os.Chdir("../../test-fixtures/userland")
	
	config, _ := habitat.GetConfig()
	loader := loaders.NewCSSFromConfig(config)

	if err := loader.Build(); err != nil { 
		t.Fatal("Build failed; ", err) 
		return
	}


	_, err := ioutil.ReadFile( loader.TargetFile )
	if err != nil { 
		t.Fatal("Could not open output file: ", err)
		return 
	}

	//TODO: Check content of file :\ 
}