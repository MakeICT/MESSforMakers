package views

import (
	//"fmt"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/makeict/MESSforMakers/models"
)

// View is a struct holding a template cache with a render method defined on it.
type View struct {
	TemplateCache map[string]*template.Template
}

// TemplateData is a holder for all the default data, and an interface for the rest
type TemplateData struct {
	AuthUser  *models.User
	CSRFToken string
	Flash     string
	Root      string
	PageTitle string
	Data      map[string]interface{}
}

// Render writes the template and data to the provided writer
func (v *View) Render(w http.ResponseWriter, r *http.Request, page string, td *TemplateData) error {
	t := v.TemplateCache[page]
	return t.ExecuteTemplate(w, page, td)

}

// LoadTemplates takes a string of folders and loads the templates into the view
func (v *View) LoadTemplates(f string) error {

	tc := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("templates/%s/*.gohtml", f))
	if err != nil {
		return fmt.Errorf("could not find view page templates: %v", err)
	}

	for _, p := range pages {
		n := filepath.Base(p)

		t, err := template.New(n).ParseFiles(p)
		if err != nil {
			return fmt.Errorf("could not create template from page: %v", err)
		}

		t, err = t.ParseGlob("templates/layouts/*.gohtml")
		if err != nil {
			return fmt.Errorf("could not create layout templates: %v", err)
		}

		t, err = t.ParseGlob(fmt.Sprintf("templates/%s/include/*.gohtml", f))
		if err != nil {
			return fmt.Errorf("could not create include templates: %v", err)
		}

		tc[n] = t
	}

	v.TemplateCache = tc

	return nil

}

//AddMap is useful for controllers when there are many values to be added to the DataStore
func (d *TemplateData) AddMap(ms ...map[string]interface{}) {
	if d.Data == nil {
		d.Data = make(map[string]interface{})
	}
	for _, stringMap := range ms {
		for k, v := range stringMap {
			d.Data[k] = v
		}
	}
}

//Add is useful for views when just a single value needs to be added
func (d *TemplateData) Add(k string, v interface{}) {
	if d.Data == nil {
		d.Data = make(map[string]interface{})
	}
	d.Data[k] = v
}
