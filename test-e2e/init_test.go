package e2e

import (
	"os"
	"testing"

	"github.com/jacobconley/habitat/habconf"
)

func TestSetup(t * testing.T) { 
	os.Chdir("../test-fixtures/userland")

	_, err := habconf.LoadConfig() 
	if err != nil { 
		t.Error("Error initializing config: ", err)
	}
}