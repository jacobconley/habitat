package e2e

import (
	habitat "habitat/src"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionOpen(t * testing.T) { 
	os.Chdir("../../test-fixtures/userland")

	config, err := habitat.GetConfig() 
	assert.Nil(t, err) 


	
	conn, err := config.NewConnection()
	assert.Nil(t, err) 

	assert.Nil(t, conn.Open())
	assert.Nil(t, conn.Close())
}