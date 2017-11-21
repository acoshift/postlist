package view

import (
	"html/template"
	"path/filepath"
)

var (
	tmplIndex  = loadTemplate("index", "layout")
	tmplCreate = loadTemplate("create", "layout")
)

const templateDir = "template"

func loadTemplate(filenames ...string) *template.Template {
	t := template.Must(template.New("").ParseFiles(joinTemplateDir(filenames...)...))
	t = t.Lookup("root")
	if t == nil {
		panic("template root not found")
	}
	return t
}

func joinTemplateDir(filenames ...string) []string {
	xs := make([]string, len(filenames))
	for i, filename := range filenames {
		xs[i] = filepath.Join(templateDir, filename+".tmpl")
	}
	return xs
}
