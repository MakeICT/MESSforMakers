package views

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// app will initialize the views and store each view in a sub struct.
// a view will be a template and
// app will have a method to addDefaultData to a viewer
// the viewer interface will require the methods setCSRFtoken, setAuthenticateduser, setFlash
// a controller then expects to have access to the app struct. Can call adddefaultdata and pass it a view

//generic defaults pages for each controller
type View struct {
	Templates *template.Template
	Name      string
}

type ViewData struct {
	//template data that applies to all pages on the site (CSRF tokens, logged in username, etc)
	CSRFToken         string
	AuthenticatedUser string
	Flash             string
}

type Viewer interface {
	SetCSRFToken(string)
	SetAuthenticatedUser(string)
	SetFlash(string)
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, layout string, td interface{}) error {
	return self.Template.ExecuteTemplate(w, self.Layout, data)
}

func (v *View) LoadTemplates() {
	//load sitewide templates and fragments
	//load view-specific templates and fragments
	sitefiles, err := filepath.Glob("templates/layouts/*.gohtml")
	if err != nil {
		log.Panic(err)
	}

	viewfiles, err := filepath.Glob(fmt.Sprintf("templates/%s/*.gohtml", v.Name))
	if err != nil {
		log.Panic(err)
	}
	files := append(sitefiles, viewfiles...)

	v.Templates = template.Must(template.New("index").ParseFiles(files...))

	return files
}
