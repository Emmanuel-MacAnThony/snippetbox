package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Emmanuel-MacAnThony/snippetbox/internal/models"
	
)

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
	CurrentYear int
	Form any
}

func humanDate(t time.Time) string{
	return t.Format("02 Jan 2006 at 15:04")
}

var templateFunctions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache()(map[string]*template.Template, error){

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")

	if err != nil{
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		// parse the base template into a template set
		ts, err := template.New(name).Funcs(templateFunctions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		
		cache[name] = ts

	}

	return cache, nil

}