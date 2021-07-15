package server

import (
	"database/sql"
	"fmt"

	"github.com/jacobconley/habitat/habconf"
	"github.com/rs/zerolog/log"
)

type ErrorHandlers struct { 

	Raw func(Context, error) (string, bool, error)
	//TODO: Other renderers https://github.com/jacobconley/habitat/issues/35

}




type HttpError struct {
	Status int 
	Message string 
} 

func (err HttpError) Error() string { 
	a := fmt.Sprintf("HTTP %d", err.Status)
	if err.Message == "" { 
		return a
	} else {
		return fmt.Sprintf("%s: %s", a, err.Message)
	}
}




func errorToStatusCode(err error) int { 

	if err == sql.ErrNoRows { 
		return 404
	}

	return 500 

}



//NEXT: Make the actual function that handles all errors, performs whatever logic to figure out what to render
func (handlers ErrorHandlers) auto(server * Server, rtype renderType, context Context, err error) { 
	log.Debug().Msg("invoking ErrorHandlers")

	rw := context.ResponseWriter
	writeCode := func(e error){ 
		rw.WriteHeader(errorToStatusCode(e))
	}

	// As we add render types, expand this logic, if the config allows fallbacks
	if handlers.Raw != nil { 
		output, handled, e2 := handlers.Raw(context, err)

		if e2 != nil { 
			log.Err(e2).Msg("occured while rendering ErrorHandlers.Raw")

			rw.WriteHeader(500)
			if habconf.Errors.FallbackToHabitatTemplate { 
				defaultErrorHandlers.Raw(context, err) 
			}

		} else if handled { 

			writeCode(err)
			rw.Write( []byte(output) )
			return 

		}
	}
	

}






var defaultErrorHandlers = ErrorHandlers { 
	Raw: func(c Context, e error) (string, bool, error) {
		return "Default error handler!", true, nil
	},
}