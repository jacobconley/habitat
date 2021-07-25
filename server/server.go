package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Server routes and serves all HTTP requests for Habitat applications.
// Where possible, its methods mimic those of Buffalo's App, but their arguments reflect the underlying parameters to Mux.

type Server struct { 
	Mux * mux.Router

	ErrorHandlers 		ErrorHandlers
}

// ServeHTTP forwards to Mux.ServeHTTP
func (r * Server) ServeHTTP(w http.ResponseWriter, req *http.Request) { 
	r.Mux.ServeHTTP(w, req)
}


func NewServer() *Server { 
	mux := mux.NewRouter()

	mux.HandleFunc("/!/assets/{path:[^:]+}:{hash:[a-fA-F0-9]+}", HandleAsset)

	srv :=  &Server{ Mux: mux }

	//TODO:  [ISSUE #35] Change these to HTML when we add it 
	mux.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, req * http.Request) { 
		ctx := srv.NewContext(rw, req) 
		srv.handleError(errNotFound, renderString, ctx)
	})
	mux.MethodNotAllowedHandler = http.HandlerFunc(func(rw http.ResponseWriter, req * http.Request) { 
		ctx := srv.NewContext(rw, req)
		srv.handleError(errMethodNotAllowed, renderString, ctx)
	})

	return srv
}



func (r * Server) NewContext(rw http.ResponseWriter, req * http.Request) * Context { 
	return &Context{
		Request: req, 
		Response: rw,
	}
}




type Matcher struct { 
	server * Server 
	path string // to HandleFunc
}



func (r * Server) Match(path string) Matcher { 
	return Matcher { 
		server: r,
		path: path,
	}
}

// TODO: MatchPrefix (Server -> Matcher)


func (m Matcher) GET() Renderer { 
	return Renderer{
		Matcher: m,
		methods: []string{"GET"},
	}
}

// TODO: All filters like Headers, Queries, Methods etc - Check Mux documentation - return another Matcher
// TODO: Other request methods - return Renderer
// TODO: Resource https://github.com/jacobconley/habitat/issues/29 - return whatever object renders a Resource, we will use reflection



// // ConfigureRoutes weird snippet if we decide to do a block initialization
// // but this should probably be up at the app level?  
// // On the other hand that's more restrictive than I'd like to be for this project
// func ConfigureRoutes( block func(r * Server) ) * Server { 
// 	r := NewServer()
// 	block(r) 
// 	return r
// 	// Should it return that?  mayb should just be a smaller interface without the route functions supposed to be called in the block
// }
/*
	One problem the above was supposed to solve was being able to detect not-allowed methods
		even when multiple allowed on the same route
	Simply reusing the intermediate config object could also solve that problem, 
	At first I was going to try to keep it declarative but
		we can just infer from hindsight since the strings will be equal 
*/


/*

	x := r.Match(route). // or MatchPrefix
		Headers(a,b).
		Queries(a,b).			<-- Matcher
		...
		Methods(...) 
		 - OR - 				<-- Renderer
		GET().HTML(...) ;
		POST().JSON(...) ; 

	If they want to handle multiple, they can just assign a variable like that


	Context:  It can carry a presumed "default" subdirectory kinda like rails does, mimicing controllers?
		has to be opt out tho
		
 */