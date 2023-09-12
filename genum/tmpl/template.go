package tmpl

import (
	_ "embed"
	"text/template"
)

//go:embed enumTemplate.gotmpl
var tmpl string

var EnumTemplate = template.Must(template.New("enum").Parse(tmpl))
