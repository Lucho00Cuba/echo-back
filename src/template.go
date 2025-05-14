package main

import (
	"html/template"
	"log"
)

// tmpl is the shared parsed HTML template.
var tmpl *template.Template

// LoadTemplate parses the HTML template at the given path.
// It is called once during initialization.
func LoadTemplate(path string) {
	var err error
	tmpl, err = template.ParseFiles(path)
	if err != nil {
		log.Fatalf("failed to load template: %v", err)
	}
}

// GetTemplate returns the current parsed template.
// Used in handlers.
func GetTemplate() *template.Template {
	return tmpl
}
