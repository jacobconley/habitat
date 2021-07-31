package server

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)



func testServer() (request * http.Request, recorder * httptest.ResponseRecorder, server * Server) { 
	req, err := http.NewRequest("GET", "/foo", nil) 
	if err != nil { 
		panic(err)
	}

	rec := httptest.NewRecorder()
	
	srv := NewServer()

	return req, rec, srv 

}


func testAssertHTML(t *testing.T, rec * httptest.ResponseRecorder) (doc * goquery.Document) { 

	assert.Equal(t, http.StatusOK, rec.Code, "Status code")
	assert.NotEqual(t, 0, rec.Body.Len(), "Nonzero body")

	
	doc, err := goquery.NewDocumentFromReader(rec.Body)
	if assert.NoError(t, err, "parsing html") { 
		return doc 
	} else { 
		return nil 
	}

}


func TestStringHTML(t *testing.T) { 
	req, rec, srv := testServer()

	srv.Match("/foo").GET().HTMLString(func(hab *Context, html *HTMLRenderer) (output template.HTML, err error) {
		return "<p>chicken nuggies</p>", nil
	})

	srv.ServeHTTP(rec, req) 
	doc := testAssertHTML(t, rec) 
	if doc != nil { 
		find := doc.Find("p")
		assert.Equal(t, 1, find.Length(), "# paragraphs")
		assert.Equal(t, "chicken nuggies", find.First().Text())
	}

}


func TestAddCSS(t *testing.T) { 
	req, rec, srv := testServer()

	srv.Match("/foo").GET().HTMLString(func(hab *Context, html *HTMLRenderer) (output template.HTML, err error) {
		html.AddCSS("styles-p.css")
		return "<p>chicken nuggies</p>", nil
	})
	
	srv.ServeHTTP(rec, req) 
	doc := testAssertHTML(t, rec) 
	if doc != nil { 
		find := doc.Find("link")
		assert.Equal(t, 1, find.Length(), "# links")
	}
	

}

func TestAddJS(t *testing.T) { 
	req, rec, srv := testServer()

	srv.Match("/foo").GET().HTMLString(func(hab *Context, html *HTMLRenderer) (output template.HTML, err error) {
		html.AddJS("main.js")
		return "<p>chicken nuggies</p>", nil
	})
	
	srv.ServeHTTP(rec, req) 
	doc := testAssertHTML(t, rec) 
	if doc != nil { 
		find := doc.Find("script")
		assert.Equal(t, 1, find.Length(), "# scripts")
	}

}