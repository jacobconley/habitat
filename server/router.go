package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router routes and serves all HTTP requests for Habitat applications.
// Where possible, its methods mimic those of Buffalo's App, but their arguments reflect the underlying parameters to Mux.

type Router struct { 
	Mux * mux.Router
}

// ServeHTTP forwards to Mux.ServeHTTP
func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) { 
	r.Mux.ServeHTTP(w, req)
}


type Matcher struct { 
	Router 
	path string // to HandleFunc
}


func NewRouter() *Router { 
	mux := mux.NewRouter()

	mux.HandleFunc("/!/assets/{path:[^:]+}:{hash:[a-fA-F0-9]+}", HandleAsset)

	return &Router{ Mux: mux }
}



func (r Router) Match(path string) Matcher { 
	return Matcher { 
		Router: r,
		path: path,
	}
}

// TODO: MatchPrefix (Router -> Matcher)


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
// func ConfigureRoutes( block func(r * Router) ) * Router { 
// 	r := NewRouter()
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