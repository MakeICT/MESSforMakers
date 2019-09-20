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
	ViewData  interface{}
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
