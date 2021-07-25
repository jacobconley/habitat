package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/jacobconley/habitat/habconf"
	"github.com/rs/zerolog/log"
)

type ErrorHandlers struct { 

	String func(* Context, error) (string, bool, error)
	//TODO: Other renderers https://github.com/jacobconley/habitat/issues/35

}




// we can wrap every error in an http error
// the user then calls .Unwrap if they need the root error 
// this also allows user to edit status, etc 

// ... or we just use normal error and add a int * status 
// cos if we do the above it'd have to be by reference 

// we also have to take into account how this interacts with 
// hab context, 
// and where the errorToStatusCode result comes into play there 

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

	if err == errMethodNotAllowed { 
		if habconf.Errors.RenderHTTPErrors.MethodNotAllowed { 
			return http.StatusMethodNotAllowed
		} else { 
			return http.StatusNotFound
		}
	}

	if err == errNotFound { 
		return http.StatusNotFound
	}

	if err == sql.ErrNoRows { 
		return http.StatusNotFound
	}

	return http.StatusInternalServerError

}



func (srv Server) handleError(err error, rtype renderType, context * Context) { 

	log.Err(err).Msg("")
	log.Debug().Msg("handleError invoked")

	handlers := srv.ErrorHandlers
	rw := context.Response

	context.Status = errorToStatusCode(err)

	// As we add render types, expand this logic, check if the config allows fallbacks, test alladat 
	if handlers.String != nil { 
		output, handled, e2 := handlers.String(context, err)

		if e2 != nil { 
			log.Err(e2).Msg("occured while rendering ErrorHandlers.String")

			rw.WriteHeader(500)
			if habconf.Errors.FallbackToHabitatTemplate { 
				out, _, _ :=defaultErrorHandlers.String(context, err) 
				context.writeOut( []byte(out) )
			}
			return 

		} else if handled { 

			context.writeOut([]byte(output))
			return 

		} else { 
			log.Debug().Msg("Error was ignored by handlers")
		}

	} 
	
	if habconf.Errors.FallbackToHabitatTemplate { 

		//TODO: [ISSUE #35] change to html when its added 
		log.Debug().Msg("falling back to default")
		out, _, _ := defaultErrorHandlers.String(context, err)
		context.writeOut( []byte(out) )
		return 

	}
}



var defaultErrorStrings = struct{
	def string
	notFound string
} { 
	def: "An unexpected error occured.",
	notFound: "Page not found",
}



var defaultErrorHandlers = ErrorHandlers { 
	String: func(c * Context, e error) (string, bool, error) {
		if e == errNotFound { 
			return defaultErrorStrings.notFound, true, nil 
		}
		return defaultErrorStrings.def, true, nil
	},
}
