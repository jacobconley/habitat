package e2e

import (
	"os"
	"testing"

	"github.com/jacobconley/habitat/habconf"
	"github.com/stretchr/testify/assert"
)

func TestConnectionOpen(t * testing.T) { 
	os.Chdir("../test-fixtures/userland")

	config, err := habconf.LoadConfig() 
	assert.Nil(t, err) 


	
	conn, err := config.NewConnection()
	assert.Nil(t, err) 

	assert.Nil(t, conn.Open())
	assert.Nil(t, conn.Close())
}