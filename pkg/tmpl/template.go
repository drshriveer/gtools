package tmpl

import (
	"embed"
	"html/template"
)

//go:embed *.gotmpl
var embededFS embed.FS

var EnumTemplate = template.Must(template.ParseFS(embededFS, "enumTemplate.gotmpl"))
