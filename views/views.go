package views

import (
	//"fmt"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/makeict/MESSforMakers/models"
)

// a view is a struct holding a template cache with a render method defined on it.
type View struct {
	TemplateCache *template.Template
}

// A templatedata is a holder for all the default data, and an interface for the rest
type TemplateData struct {
	AuthUser  *models.User
	CSRFToken string
	Flash     string
	Root      string
	PageTitle string
	Data      map[string]interface{}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, layout string, td *TemplateData) error {

	return v.TemplateCache.ExecuteTemplate(w, layout, td)

}

func (v *View) LoadTemplates(ff []string) error {
	var files []string
	for _, f := range ff {
		fg, err := filepath.Glob(fmt.Sprintf("templates/%s/*.gohtml", f))
		if err != nil {
			return fmt.Errorf("could not find template files: %v", err)
		}
		files = append(files, fg...)
	}
	//fmt.Println(files)
	tc, err := template.ParseFiles(files...)
	if err != nil {
		return err
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
