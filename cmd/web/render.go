package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

var pathToTemplates = "./cmd/web/templates"

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	DataMap       map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	// User            *data.user
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", pathToTemplates),
		fmt.Sprintf("%s/header.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/navbar.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/footer.partial.gohtml", pathToTemplates),
		fmt.Sprintf("%s/alerts.partial.gohtml", pathToTemplates),
	}
	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", pathToTemplates, t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	if td == nil {
		td = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.ErrorLog.Println("error parsing template files", err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err := tmpl.Execute(w, app.AddDefaultData(td, r)); err != nil {
		app.ErrorLog.Println("error executing template files", err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func (app *Config) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	if app.isAuthenticated(r) {
		td.Authenticated = true
		// TODO: get more user information
	} else {
		td.Authenticated = false
	}
	td.Now = time.Now()
	return td
}

func (app *Config) isAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
