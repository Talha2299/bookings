package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Talha2299/bookings/internal/config"
	"github.com/Talha2299/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var pathToTemplates = "./templates"

var functions = template.FuncMap{
	"humanDate": HumanDate,
	"formatDate": FormatDate,
	"iterate": Iterate,
	"add": Add,
}

var app *config.AppConfig

// NewRenderer set the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func Add(a, b int) int {
	return a + b
}
// Iterate return a slice of ints, starts at 1 and going to count
func Iterate(count int) []int {
	var i int
	var items []int

	for i=0;i<count;i++ {
		items = append(items, i)
	}
	return items
}

// HumanDate reutrn time in YYYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// AddDefaultData adds data for all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// Template for render templates
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		//get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// this is just use for testing thats why we rebuilt
		// the cache on every request
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		return errors.New("couldnt get template from template caches")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		fmt.Println("error  writng templates to browser", err)
		return err
	}

	return nil
}

//createTemplateCache create template cahce as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, nil
}
