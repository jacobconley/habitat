package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func testStatusBody(t *testing.T, srv * Server, method string, path string, status int, body string) { 
	req, err := http.NewRequest(method, path, nil)
	assert.NoError(t, err, "initializing request")

	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	assert.Equal(t, status, rec.Code, "Status code")
	assert.Equal(t, body, rec.Body.String(), "Response body")
}

func TestDefault404(t *testing.T) {
	testStatusBody(t, NewServer(), 
		"GET", "/asdf-reee", 
		404, defaultErrorStrings.notFound)
}

func TestCustom404(t *testing.T) { 
	srv := NewServer()
	strang := "Caint find it my guy"

	srv.ErrorHandlers.Raw =  func(hab * Context, err error) ( string, bool, error ){ 
		return strang, true, nil
	}

	testStatusBody(t, srv,
		"GET", "/asdf-ree",
		404, strang)
}