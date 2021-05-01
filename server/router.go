package server

import "github.com/gorilla/mux"

/*
HabitatRouter returns the primary router for serving Habitat
	TODO: Include plugins
*/
func Router() *mux.Router { 
	r := mux.NewRouter()

	r.HandleFunc("/!/assets/{path:[^:]+}:{hash:[a-fA-F0-9]+}", HandleAsset)

	return r
}