package server

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestUnmarshalURLVals(t *testing.T) { 

	float := 1.234

	out := struct{
		Foo string 
		Integer int 
		Float float64
	}{}

	vals := url.Values{
		"foo": 			[]string{ "test" },
		"integer": 		[]string{ "23" },
		"float": 		[]string{ fmt.Sprintf("%f", float) },
	}


	err := unmarshalURLValues(vals, &out)
	assert.NoError(t, err)

	assert.Equal(t, out.Foo, "test") 
	assert.Equal(t, out.Integer, 23)
	assert.Equal(t, out.Float, float)

}


func TestUnmarshalURLValsNameCase(t *testing.T) { 
	out := struct {
		FooBar 	string 
		AInt 	int
		Foobar 	string 
		BInt 	int
	} { } 

	vals := url.Values{
		"Foobar": 	[]string{ "asdf" },
		"bint": 	[]string{ "23" }, 
		"aint": 	[]string{ "12" },
	}

	err := unmarshalURLValues(vals, &out) 
	assert.NoError(t, err)

	assert.Equal(t, out.Foobar, "asdf")
}

func TestUnmarshalURLValsNameTags(t *testing.T) { 
	out := struct{
		Foo 	string 	`url:"name"`
		Bar 	string 	`json:"foo"`
		Omit 	string 	`json:"omit,omitempty"`
		Baz 	int
	} {}

	vals := url.Values { 
		"foo": 	[]string{ "123" },
		"bar": 	[]string{ "asdf" },
		"baz": 	[]string{ "890" },
		"name": []string{ "jacob" },
		"omit": []string{ "aaa" },
	}

	err := unmarshalURLValues(vals, &out) 
	assert.NoError(t, err)

	assert.Equal(t, out.Foo, "jacob")
	assert.Equal(t, out.Bar, "123")
	assert.Equal(t, out.Baz, 890) 
	assert.Equal(t, out.Omit, "aaa")
}

func TestUnmarshalURLValsIgnore(t *testing.T) { 
	out := struct { 
		Foo 	string
		Bar 	string 	`url:"-"`
		Baz 	string 	`json:"-"`
	}{} 

	vals := url.Values { 
		"foo": 	[]string{ "ree" },
		"bar": 	[]string{ "asdf" },
		"baz": 	[]string{ "jkl;" },
	}

	err := unmarshalURLValues(vals, &out) 
	assert.NoError(t, err)

	assert.Equal(t, out.Foo, "ree")
	assert.Equal(t, out.Bar, "")
	assert.Equal(t, out.Baz, "")
}


func TestUnmarshalURLValsErrParseFloat(t *testing.T) { 
	out := struct { 
		Foo  float64 
	} {}

	vals := url.Values { 
		"foo": 	[]string{ "asdf" },
	}

	err := unmarshalURLValues(vals, &out)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errRequestParam)
}




func TestUnmarshalURLValsNilPointers(t * testing.T) { 
	out := struct { 
		Thang 	int 
		Guy 	*int 
		ZeroGuy	int 
	}{}

	vals := url.Values { 
		"thang": []string{ "22" },
	}

	err := unmarshalURLValues(vals, &out)
	assert.NoError(t, err) 

	assert.Equal(t, out.Thang, 22) 
	assert.Nil(t, out.Guy)
	assert.Equal(t, out.ZeroGuy, 0) 
}


func TestUnmarshalURLValsEmbeddedStructs(t * testing.T) { 
	type Embed struct { 
		Bar string 
	}
	out := struct { 
		Foo string 
		Embed
	}{}

	vals := url.Values{ 
		"Foo": []string{ "ree" },
		"Bar": []string{ "asdf" },
	}

	err := unmarshalURLValues(vals, &out) 
	assert.NoError(t, err)

	assert.Equal(t, out.Foo, "ree")
	assert.Equal(t, out.Bar, "asdf")
}


func TestUnmarshalURLValsSlices(t * testing.T) { 

	out := struct { 
		Single 	string
		Strings []string
		Ints 	[]int
	}{} 

	vals := url.Values { 
		"single": 	[]string{ "abc", "def" },
		"strings":	[]string{ "foo", "bar" },
		"ints": 	[]string{ "123", "456" },
	}

	err := unmarshalURLValues(vals, &out) 
	assert.NoError(t, err)

	assert.Equal(t, out.Single, "abc")
	assert.ElementsMatch(t, out.Strings, []string{ "foo", "bar" })
	assert.ElementsMatch(t, out.Ints, []int{ 123, 456 })
}




func TestUnmarshalURLValsRequired(t * testing.T) { 
	out := struct { 
		Present 	int 
		Missing 	string `http:"required"`
	}{}  

	vals := url.Values{ 
		"Present": []string{ "22" },
	}

	err := unmarshalURLValues(vals, &out) 
	assert.ErrorIs(t, err, errRequestParam)
}


func TestUnmarshalURLValsNonempty(t * testing.T) { 
	out := struct { 
		Present 	int 
		Empty 		string `http:"nonempty"`
	}{}  

	vals := url.Values{ 
		"Present": []string{ "22" },
		"Empty": []string{ "" },
	}

	err := unmarshalURLValues(vals, &out) 
	assert.ErrorIs(t, err, errRequestParam)
}