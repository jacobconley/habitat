package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"

	"github.com/jacobconley/habitat/habconf"
)

// ======
// Errors
// ------
// Be sure to keep these errors up-to-date with errorToStatusCode in error-handling.go
// 	- or isErrBadRequest below
// ------


var ErrNoBody 				= errors.New("Missing request body")
var ErrNoContentType 		= errors.New("Missing Content-Type header")
var	ErrBodyTooLarge 		= errors.New("Request body too large")

var ErrBadContentType 		= errors.New("Unsupported Content-Type")
var ErrExpectedForm 		= errors.New("ParamValues() expected Content-Type: application/x-www-form-urlencoded - use ParamUnmarshal() for other content types")

var errRequestParam 		= errors.New("Missing or invalid request parameter")

// ErrParamMisisng represents an error caused by a required parameter that is not supplied
type ErrParamMissing struct { 
	// Name is the name of the parameter 
	Name string 
	PresentButEmpty bool 
}
func (err ErrParamMissing) Error() string { 
	if err.PresentButEmpty { 
		return fmt.Sprintf("Parameter '%s' is empty", err.Name)
	} else { 
		return fmt.Sprintf("Missing parameter '%s'", err.Name)
	}
}
func (err ErrParamMissing) Is(e2 error) bool { 
	// signal to error type helpers below
	return e2 == errRequestParam
}


// ErrParamParse represents a wrapped error caught while attempting to parse a given parameter 
type ErrParamParse struct { 
	// Name is the name of the parameter 
	Name string 

	ExpectedType string 
	Cause error
}
func (err ErrParamParse) Error() string { 
	return fmt.Sprintf("Could not parse parameter '%s' as type '%s'", err.Name, err.ExpectedType)
}
func (err ErrParamParse) Is(e2 error) bool { 
	// signal to error type helpers below
	return e2 == errRequestParam
}
func (err ErrParamParse) Unwrap() error { 
	return err.Cause
}



type ErrParamUnsupportedType struct { 
	Name 		string 
	FieldKind 	string 
}
func (err ErrParamUnsupportedType) Is(e2 error) bool { 
	return e2 == errRequestParam
}
func (err ErrParamUnsupportedType) Error() string { 
	return fmt.Sprintf("Could not parse parameter '%s'; unsupported kind of field '%s'", err.Name, err.FieldKind)
}





// 
// Error-type helpers, used to figure out what HTTP codes to return
//


// isErrBadRequest returns true if the given error is one of the request handling errors in habitat/server, and therefore should return HTTP 400 
func isErrBadRequest(err error) bool { 
	return err == ErrNoBody || err == ErrNoContentType || errors.Is(err, errRequestParam)
}

// isErrUnprocessableEntity returns true only if the given error is a server.ErrParam* error AND habconf.Errors.RenderHTTPErrors.UnprocessableEntity is true
func isErrUnprocessableEntity(err error) bool { 
	return habconf.Errors.RenderHTTPErrors.UnprocessableEntity && errors.Is(err, errRequestParam)
}






// readBody reads the entire request body, limited by habconf.MaxFormSize, and closes r.Body 
func readBody(r * http.Request) (res []byte, err error) { 	

	defer r.Body.Close()

	buf := make([]byte, 8<<3)
	var size int64

	var n int
	for { 
		n, err = r.Body.Read(buf)

		size = size + int64(n) 
		if size > habconf.MaxFormSize { 
			err = ErrBodyTooLarge 
			return 
		}

		if n > 0 { 
			res = append(res, buf[:n]...)
		}
		if err == io.EOF { 
			err = nil 
			return 
		}
		if err != nil { 
			return 
		}
	}
}


// getContentType gets the header and runs it through mime.ParseMediaType, or returns ErrNoContentType
func getContentType(r * http.Request) (ct string, err error) { 
	ct = r.Header.Get("Content-Type")
	if ct == "" {
		err = ErrNoContentType
		return 
	}
	ct, _, err = mime.ParseMediaType(ct)
	return
}


// We ignore the possibility of GET requests having message bodies - 
// see here https://stackoverflow.com/a/983458


func reqIsJSON(r * http.Request) (bool, error) { 
	if r.Method == "GET" { 
		return false, nil 
	}

	ct, err := getContentType(r) 
	if err != nil { 
		return false, err 
	}

	return ct == "application/json"|| ct == "text/x-json", nil
}



// RequestValues returns the `url.Values` representing this request's parameters.
// If the request is a GET, this function returns URL.Query().  
// If it is a POST, the body will be parsed as a query string if the Content-Type is application/x-www-form-urlencoded, otherwise ErrExpectedForm will be returned. 
func (hab Context) RequestValues() (vs url.Values, err error) { 
	r := hab.Request 

	if r.Method == "GET" { 
		return r.URL.Query(), nil 
	} 


	var ctype string 
	ctype, err = getContentType(hab.Request)
	if err != nil { 
		return  
	}
	if r.Body == nil { 
		err = ErrNoBody 
		return 
	}

	if ctype != "application/x-www-form-urlencoded" { 
		var body []byte 
		body, err = readBody(r)
		if err != nil { 
			return 
		}

		return url.ParseQuery(string(body))

	} else { 
		return nil, ErrExpectedForm
	}
}


// RequestUnmarshal attempts to unmarshal the request's parameters into the given argument.
// If the request is a JSON request (non-GET with JSON Content-Type headers), json.Unmarshal will be invoked.
// Otherwise, this function will attempt to unmarshal as a url-encoded string in the URL query or request body if applicable. 
func (hab Context) RequestUnmarshal(into interface{}) error { 
	r := hab.Request


	body, err := readBody(r) 
	if err != nil { 
		return err 
	}


	isJSON, err := reqIsJSON(r)
	if err != nil { 
		return err 
	}


	if isJSON { 
		return json.Unmarshal(body, into)
	} else { 
		vs, err := hab.RequestValues()
		if err != nil { 
			return err
		}
		return unmarshalURLValues(vs, into)
	}


}


