package server

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// reflectForEachField invokes `op` for each reflect field.  `name` is the name of the parameter determined by the tags or field name.
func reflectForEachField(rtype reflect.Type, rval reflect.Value, op func(name string, field reflect.StructField, fieldVal reflect.Value)error ) (err error) { 

	for i := 0; i < rval.NumField(); i += 1 { 
		fieldVal  	:= rval.Field(i) 
		field 		:= rtype.Field(i)


		// Embedded structs 
		if field.Anonymous && field.Type.Kind() == reflect.Struct { 
			err = reflectForEachField(field.Type, fieldVal, op)
			if err != nil { 
				return 
			}
		}


		var name string
		if attr, ok := field.Tag.Lookup("json"); ok { 
			name = strings.Split(attr, ",")[0]
		} else if attr, ok := field.Tag.Lookup("url"); ok { 
			name = attr 
		} else { 
			name = field.Name
		}

		if name == "-" { 
			return 
		}

		err = op(name, field, fieldVal) 
	}

	return
}



func reflectSetURLValues(rtype reflect.Type, rval reflect.Value, vs url.Values) error { 

	return reflectForEachField(rtype, rval, func(name string, field reflect.StructField, fieldVal reflect.Value) error { 


		// like json, prefer exact-case match but accept case-insensitive match 
		// https://pkg.go.dev/encoding/json#Unmarshal
		// Though this technically isn't quite in line with URL spec - https://stackoverflow.com/a/24700171
		// But exported field names have to be uppercase which looks ugly and unusual in URLs; and supposedly Microsoft does it anyways
		// Should we signal an inexact match somehow?  Perhaps if a use case comes up for it... 
		// Prefer tag to field match too?  we'd have to rethink the structure of reflectForEachField

		var values []string 
		hasValue := false 

		for urlKey, urlVal := range vs { 

			if name == urlKey { 
				// prefer exact match
				values = urlVal
				hasValue = true 
				break 
			}

			if strings.EqualFold(name, urlKey) { 
				// this will only be executed if no exact match has been found, including on the current key/name
				values = urlVal 
				hasValue = true 
			}
		}

		tag, hasTag := field.Tag.Lookup("http")

		if !hasValue { 
			if hasTag && tag == "required" { 
				return ErrParamMissing {
					Name: name,
					PresentButEmpty: false,
				}
			} else { 
				return nil 
			}
		}


		// If this tag is present, it will verify that the provided value is nonzero
		// Except for booleans, which semantically are only ever zero or nonzero 
		nonempty := tag == "nonempty"
		

		// Supporting arrays, and nullable pointers like JSON 
		fieldKind := field.Type.Kind()
		var targetKind reflect.Kind
		var targetVal reflect.Value
		variadic := false // this will make supporting Arrays easier later; see below 

		if fieldKind == reflect.Ptr { 
			targetKind 		= field.Type.Elem().Kind()
			targetVal	 	= fieldVal.Elem()
		} else if fieldKind == reflect.Slice {
			targetKind = field.Type.Elem().Kind()
			targetVal = fieldVal
			variadic = true 
		} else if fieldKind == reflect.Array { 
			// Not yet supported; not sure yet what the implications of this are for reflection
			// https://github.com/jacobconley/habitat/issues/42
			return ErrParamUnsupportedType{
				Name: field.Name,
				FieldKind: fmt.Sprintf("Array of %s", field.Type.Elem().Kind().String()),
			}
		} else { 
			targetKind 		= fieldKind
			targetVal 		= fieldVal
		}




		if variadic { 
			len := len(values)
			slice := reflect.MakeSlice(field.Type, len, len)
			
			for i, input := range(values) { 
				parsed, err := stringToValue(field.Name, input, targetKind, nonempty)
				// Not exactly sure how to handle `nonempty` for slices - could be some funny semantics afoot past the first value? 
				if err != nil { 
					return err 
				}

				slice.Index(i).Set( reflect.ValueOf(parsed) )
			}

			targetVal.Set(slice)
			return nil

		} else { 

			input := values[0] // I'm fairly certain we can assume there's at least one value if it makes it into the original `vs` map? 
			output, err := stringToValue(field.Name, input, targetKind, nonempty)

			if err != nil { 
				return err
			}

			outval := reflect.ValueOf(output) 
			targetVal.Set( outval )
			return nil 
		}

	})

}


// stringToValue parses a string `input` according to `targetKind` and returns an appropriate `interface{}` to be initialized as a reflect value 
func stringToValue(varname string, input string, targetKind reflect.Kind, nonempty bool) (interface{}, error) { 

	wrapParseErr := func(expectedType string, err error) error { 
		return ErrParamParse{
			Name: varname,
			ExpectedType: expectedType,
			Cause: err,
		}
	}

	emptyErr := func() error { 
		return ErrParamMissing{
			Name: varname,
			PresentButEmpty: true,
		}
	}


	switch targetKind { 

	case reflect.String:
		if nonempty && input == "" { 
			return nil, emptyErr()
		}
		return input, nil 

	case reflect.Bool:
		// From url.Values: "A setting without an equals sign is interpreted as a key set to an empty value"
		//TODO: Document this 
		return (input != "false" && input != "0"), nil

	case reflect.Float64, reflect.Float32:

		bits := 64
		if targetKind == reflect.Float32 { 
			bits = 32
		}

		out, err := strconv.ParseFloat(input, bits)
		if err != nil { 
			return nil, wrapParseErr("float", err)
		} else if out == 0.0 { 
			return nil, emptyErr()
		} else { 
			return out, nil 
		}

	case reflect.Int:
		out, err := strconv.Atoi(input)
		if err != nil { 
			return nil, wrapParseErr("int", err)
		} else if out == 0 { 
			return nil, emptyErr()
		} else { 
			return out, nil
		}



	default:
		return nil, ErrParamUnsupportedType{
			Name: varname,
			FieldKind: targetKind.String(),
		}

	}

}




func unmarshalURLValues(vs url.Values, into interface{}) error { 
	ptr_val := reflect.ValueOf(into)
	if ptr_val.Kind() != reflect.Ptr { 
		panic("Must unmarshal into a pointer")
	}

	rval := ptr_val.Elem()
	rtype := rval.Type()


	return reflectSetURLValues(rtype, rval, vs)
}