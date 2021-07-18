package server

import (
	"errors"
	"net/http"
)


type Renderer struct {
	Matcher
	methods []string 
}

func (r Renderer) allowsMethod(method string) bool { 
	if r.methods == nil { 
		return true 
		// Temporary.  If no methods array, the router only gets here by matching method
		// ^ Though to make the above situation happen, we use the .GET() postfix on the HandleFunc call I think. 
		// 	so maybe we just pass that to the user?  There's other API stuff that Mux exposes too
		// https://github.com/gorilla/mux#matching-routes
	}

	if method == "" { 
		method = "GET" // See documentation on http.Request.Method
	}

	for _, m := range r.methods { 
		if m == method { 
			return true
		}
	}
	return false
}


var errMethodNotAllowed = errors.New("HTTP Method not allowed")
var errNotFound = errors.New("HTTP Not found")


func (r Renderer) before(rw http.ResponseWriter, req * http.Request) error { 

	if !(r.allowsMethod( req.Method )) { 
		return errMethodNotAllowed
	}

	return nil 
}



// Render type definitions
// https://github.com/jacobconley/habitat/issues/35

type renderType int 
const ( 
	renderRaw renderType = iota
)




func (r Renderer) Raw( handler func(* Context) (result string, err error) ) { 

	r.server.Mux.HandleFunc( r.path, func(rw http.ResponseWriter, req *http.Request) {
 
		habctx := r.server.NewContext(rw, req)
		var res string

		err := r.before(rw, req)
		if err == nil { 
			res, err = handler(habctx) 
		}

		if err != nil { 
			r.server.handleError(err, renderRaw, habctx) 
			return
		}

		habctx.writeOut( []byte(res) )

	}) 
}


// func (r Renderer) JSON( handler func(hab * Context) (result interface{}, err error) ) { 
// }

// func (r Renderer) WebTemplate( template string, handler func(hab * Context) (vars interface{}, err error) ) { 
// }