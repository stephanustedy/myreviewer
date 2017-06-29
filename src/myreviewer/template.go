package myreviewer

import(
	"html/template"
)

var templates *template.Template

func init() {
	templates = template.Must(template.New("").ParseGlob("files/html/*.html"))
}