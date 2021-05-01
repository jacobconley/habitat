package e2e

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jacobconley/habitat/auxiliary"
)

func TestCSS(t *testing.T) {
	os.Chdir("../test-fixtures/userland")
	
	err, loader := auxiliary.BuildCSS()

	if err := loader[0].Build(); err != nil { 
		t.Fatal("Build failed") 
		return
	}


	_, err = ioutil.ReadFile( loader[0].TargetFile )
	if err != nil { 
		t.Fatal("Could not open output file: ", err)
		return 
	}

	//TODO: Check content of file :\ 
}



func TestWebpack(t *testing.T) { 
	os.Chdir("../test-fixtures/userland")

	err := auxiliary.BuildWebpack()

	if err != nil { 
		t.Fatal("Webpack build failed")
	}

}