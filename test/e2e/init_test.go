package e2e

import (
	habitat "habitat/src"
	"os"
	"testing"
)

func TestSetup(t * testing.T) { 
	os.Chdir("../../test-fixtures/userland")

	_, err := habitat.GetConfig() 
	if err != nil { 
		t.Error("Error initializing config: ", err)
	}
}