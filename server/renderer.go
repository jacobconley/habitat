package server

import (
	"net/http"

	"github.com/jacobconley/habitat/habconf"
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



func (r Renderer) before(rw http.ResponseWriter, req * http.Request) { 

	if !(r.allowsMethod( req.Method )) { 

		if habconf.GetConfig().RenderMethodNotAllowed { 
			rw.WriteHeader( http.StatusMethodNotAllowed )
			//TODO: How to render this https://github.com/jacobconley/habitat/issues/33
		}
	}
}


func (r Renderer) Raw( handler func(* Context, * http.Request) (result string, err error) ) { 

	r.Mux.HandleFunc( r.path, func(rw http.ResponseWriter, req *http.Request) {

		r.before(rw, req) 

		//TODO: Error handling https://github.com/jacobconley/habitat/issues/33
		// Above and below 

		//TODO: Provide Habitat Context to the below

		res, err := handler(nil, req) 

		if err != nil { 
			rw.WriteHeader(500) 
			return
		}

		rw.Write( []byte(res) )

	}) 
}

// func (r Renderer) JSON( handler func(hab * Context) (result interface{}, err error) ) { 
// }

// func (r Renderer) WebTemplate( template string, handler func(hab * Context) (vars interface{}, err error) ) { 
// }
//TODO: Templater type that initializes with a layout or whatever else, then can be passeed here