package main

import (
	"embed"
	"html/template"
)

//go:embed "templates"
var templateFS embed.FS

func (app *application) Documentation() (*template.Template, error) {
	tmpl, err := template.New("doc").ParseFS(templateFS, "templates/index.tmpl")
	if err != nil {
		return nil, err
	}

	return tmpl, err
}