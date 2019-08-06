package views

import (
	//"fmt"
	"html/template"
	"net/http"

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
	ViewData  interface{}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, layout string, td *TemplateData) error {

	// for _, _ := range v.TemplateCache.Templates() {
	// 	fmt.Println(t.Name())
	// }

	return v.TemplateCache.ExecuteTemplate(w, layout, td)
}
