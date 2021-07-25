package server

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)



func TestStringHTML(t *testing.T) { 
	req, err := http.NewRequest("GET", "/foo", nil) 
	assert.NoError(t, err, "initializing request") 

	rec := httptest.NewRecorder()
	
	srv := NewServer()
	srv.Match("/foo").GET().HTMLString(func(hab *Context, html *HTMLRenderer) (output template.HTML, err error) {
		return "<p>chicken nuggies</p>", nil
	})

	srv.ServeHTTP(rec, req) 

	assert.Equal(t, http.StatusOK, rec.Code, "Status code")
	assert.NotEqual(t, 0, rec.Body.Len(), "Nonzero body")

	
	doc, err := goquery.NewDocumentFromReader(rec.Body)
	assert.NoError(t, err, "parsing html")

	find := doc.Find("p")
	assert.Equal(t, 1, find.Length(), "# paragraphs")
	assert.Equal(t, "chicken nuggies", find.First().Text())
}