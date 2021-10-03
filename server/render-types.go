package server

import (
	"fmt"
	"html/template"
	"net/http"
)

// Render type definitions
// https://github.com/jacobconley/habitat/issues/35

type renderType int 
const ( 
	renderString 	renderType = iota
	renderHTMLString
	renderHTML
)


type HandlerRaw 			func(hab * Context) (err error) 
type HandlerString 			func(hab * Context) (result string, err error) 
type HandlerData 			func(hab * Context) (data interface{}, err error) 
type HandlerHTMLString		func(hab * Context, html * HTMLRenderer) (output template.HTML, err error)
type HandlerHTML			func(hab * Context, html * HTMLRenderer) (template template.Template, data interface{}, err error) 


func (r Renderer) String( action HandlerString ) { 

	r.server.Mux.HandleFunc( r.path, func(rw http.ResponseWriter, req *http.Request) {
		rtype := renderString
 
		habctx := NewContext(rw, req)
		var res string

		err := r.beforeRender(rw, req)
		if err == nil { 
			res, err = action(habctx) 
		}

		if err != nil { 
			r.server.handleError(err, rtype, habctx) 
		} else { 
			err := habctx.writeOut( []byte(res) )
			if err != nil { 
				r.server.handleError(err, rtype, habctx)
			}
		}


	}) 

}



func (r Renderer) HTMLString( action HandlerHTMLString ) { 

	r.server.Mux.HandleFunc( r.path, func(rw http.ResponseWriter, req *http.Request) { 
		rtype := renderHTMLString

		habctx := NewContext(rw, req) 
		html := new(HTMLRenderer)

		var body template.HTML
		err := r.beforeRender(rw, req)
		if err == nil { 
			body, err = action(habctx, html)
		}

		if err != nil { 
			r.server.handleError(err, rtype, habctx)
		} else { 

			html.private.OutputType 	= rtype 
			html.private.OutputString	= body 

			habctx.beforeWrite()
			err = defaultLayout.Execute(rw, htmlOutput { 
				HTMLRenderer: *html,
				htmlPrivate: html.private,
			})

			if err != nil { 
				r.server.handleError( fmt.Errorf("rendering HTML layout: %w", err), rtype, habctx) 
			}

		}

	})

}






// func (r Renderer) JSON( action func(hab * Context) (result interface{}, err error) ) { 
// }

// func (r Renderer) WebTemplate( template string, action func(hab * Context) (vars interface{}, err error) ) { 
// }