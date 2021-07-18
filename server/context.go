package server

import (
	"net/http"
)

type Context struct {
	Request * http.Request
	ResponseWriter http.ResponseWriter


	// Status is the HTTP status
	// The header can only be sent once, but this variable can be updated throughout the life cycle 
	// It can also be left unset - `http.ResponseWriter` will send a 200 by default, and Habitat will handle errors unless configured not to
	Status int
}

func (hab Context) writeOut(out []byte) { 
	if hab.Status != 0 { 
		hab.ResponseWriter.WriteHeader( hab.Status ) 
	}
	hab.ResponseWriter.Write(out)
}