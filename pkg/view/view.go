package view

import (
	"html/template"
	"net/http"
)

func render(w http.ResponseWriter, t *template.Template, data interface{}) {
	t.Execute(w, data)
}

// Index renders index view
func Index(w http.ResponseWriter, data interface{}) {
	render(w, tmplIndex, data)
}

// Create renders create view
func Create(w http.ResponseWriter) {
	render(w, tmplCreate, nil)
}
