package server

import (
	"fmt"
	"html/template"
)


type HTMLRenderer struct {
	
	LinkedCSS []string

	//TODO: Default title (simple string in server I guess)
	Title string


	private htmlPrivate

}

type htmlPrivate struct { 
	OutputType 		renderType
	OutputString 	template.HTML
	OutputTemplate 	template.Template
	OutputData 		interface{}
}




func (html * HTMLRenderer) AddCSS(path string) { 
	html.LinkedCSS = append(html.LinkedCSS, path)
}




// htmlOutput is the object upon which the template is executed. 
/*
	HTMLRenderer is an exported API.
	All data meant to be accessed from the template  has to be stored in exported variables,
		even if the data is meant to be encapsulated and unexported.
	So we keep encapsulated data in the unexported htmlData structure,
		and embed both structures in the htmlOutput structure here,
		on which we can also define whatever API functions we want accessible from the template 
 */
type htmlOutput struct { 
	HTMLRenderer
	htmlPrivate
}


func (html htmlOutput) IsTemplated() bool { 
	switch html.OutputType { 
	case renderHTML:
		return true 
	case renderHTMLString:
		return false

	default:
		panic("Unhandled HTML render type in template logic")
	}
}





var defaultLayout * template.Template

var defaultLayoutTemplate = `<html>
<head>
	<title>{{ .Title }}</title>
</head>
<body>
{{ if .IsTemplated }}
{{ template "body" .OutputData }}
{{ else }}
{{ .OutputString }}
{{ end }}
</body>
</html>`

// if we do anything else in here, move it somewhere better
// just tryign to keep all the default template logic here 
func init() { 
	t := template.New("habitat_layout")

	t , err :=  t.Parse(defaultLayoutTemplate)
	if err != nil { 
		panic(fmt.Errorf("Habitat: Could not parse default HTML layout - %w", err)) 
	}
	defaultLayout = t 

	bod := template.New("body")
	bod, err = bod.Parse("asdfjkl;")
	if err != nil { 
		panic(err)
	}
	defaultLayout.AddParseTree("body", bod.Tree)

	// HOW to clean this up 
}
