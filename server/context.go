package server

import (
	"net/http"
)

type Context struct {
	Request * http.Request

	Response http.ResponseWriter

	// Status is the HTTP status
	// The header can only be sent once, but this variable can be updated throughout the life cycle 
	// It can also be left unset - `http.ResponseWriter` will send a 200 by default, and Habitat will handle errors unless configured not to
	Status int
}


func (hab Context) beforeWrite() { 
	if hab.Status != 0 { 
		hab.Response.WriteHeader( hab.Status ) 
	}
}

func (hab Context) writeOut(out []byte) error { 
	hab.beforeWrite()
	_, err := hab.Response.Write(out)
	return err 
}


