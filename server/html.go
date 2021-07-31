package server

import (
	"fmt"
	"html/template"
)


type HTMLRenderer struct {
	
	Stylesheets 	[]Stylesheet
	Scripts 		[]Script

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




func tagAppendString(input string, key string, val string) string { 
	if val == "" { 
		return input
	} else { 
		return fmt.Sprintf(`%s %s="%s"`, input, key, val)
	}
}
func tagAppendBool(input string, key string, val bool) string { 
	if val { 
		return fmt.Sprintf("%s %s", input, key) 
	} else { 
		return input 
	}
}

type Stylesheet struct { 
	Type string 
	Href string 
}
func (s Stylesheet) String() string { 
	return fmt.Sprintf(`<link rel="stylesheet" %s />`, tagAppendString(fmt.Sprintf(`href="%s"`, s.Href), "type", s.Type))
}
func (s Stylesheet) HTML() template.HTML { 
	return template.HTML(s.String())
}


type Script struct { 

	// Attributes, referencing https://developer.mozilla.org/en-US/docs/Web/HTML/Element/script 
	Async 				bool 
	CrossOrigin 		bool 
	Defer 				bool 
	NoModule 			bool 

	Nonce 				string 
	ReferrerPolicy 		string 

	Src 				string 
	Type 				string 

	Contents 			string 

	//TODO: Integrity https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
}
func (s Script) String() string { 
	var attrStr, attrBool string 

	attrStr = tagAppendString(attrStr, "nonce", s.Nonce)
	attrStr = tagAppendString(attrStr, "referrerpolicy", s.ReferrerPolicy)
	attrStr = tagAppendString(attrStr, "type", s.Type)

	attrBool = tagAppendBool(attrBool, "async", s.Async)
	attrBool = tagAppendBool(attrBool, "crossorigin", s.CrossOrigin)
	attrBool = tagAppendBool(attrBool, "defer", s.Defer)
	attrBool = tagAppendBool(attrBool, "nomodule", s.NoModule)

	var attrs string 
	if attrStr != "" { 
		attrs = attrs + " "
	}
	attrs = attrs + attrStr
	if attrBool != "" { 
		attrs = attrs + " " 
	}
	attrs = attrs + attrBool



	return fmt.Sprintf(`<script src="%s"%s>%s</script>`, s.Src, attrs, s.Contents)
}
func (s Script) HTML() template.HTML { 
	return template.HTML(s.String())
}



func (html * HTMLRenderer) AddCSS(path string) { 
	html.Stylesheets = append(html.Stylesheets, Stylesheet{
		Type: "text/css",
		Href: path,
	})
}
func (html * HTMLRenderer) AddJS(path string) { 
	html.Scripts = append(html.Scripts, Script{
		// From MDN: "The HTML5 specification urges authors to omit the attribute rather than provide a redundant MIME type"
		Src: path,
	})
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

// TODO head template
// I don't think we'll be "replacing" defaultLayout in the rails sense
// instead we allow them to append to head, and maybe we let them provide some sort of wrapper that wraps template "body" naaaamean 

var defaultLayoutTemplate = `<html>
<head>
	<title>{{ .Title }}</title>

	{{ range .Stylesheets }}{{ .HTML }}
	{{ end }}{{ range .Scripts }}{{ .HTML }}
	{{ end }}

	{{ template "head" }}
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
	bod, err = bod.Parse("")
	if err != nil { 
		panic(err)
	}
	defaultLayout.AddParseTree("body", bod.Tree)

	head := template.New("head") 
	head, err = head.Parse("")
	if err != nil { 
		panic(err) 
	}
	defaultLayout.AddParseTree("head", head.Tree)
}
