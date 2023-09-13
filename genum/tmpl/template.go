package tmpl

import (
	_ "embed"
	"text/template"
)

//go:embed enumTemplate.gotmpl
var tmpl string

// EnumTemplate is the base template for an enum.
var EnumTemplate = template.Must(template.New("enum").Parse(tmpl))
