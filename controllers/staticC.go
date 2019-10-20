package controllers

import (
	"fmt"
	"net/http"

	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

//StaticController implements the handlers required for basic navigation pages
type StaticController struct {
	Controller
	StaticView views.View
}

//Initialize performs the required setup for a static controller
func (sc *StaticController) Initialize(cfg *util.Config, um Users, l *util.Logger) error {
	sc.setup(cfg, um, l)
	sc.StaticView = views.View{}

	if err := sc.StaticView.LoadTemplates("static"); err != nil {
		return fmt.Errorf("Error loading static templates: %v", err)
	}

	return nil
}

// Root displays the home page
func (sc *StaticController) Root() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		td, err := sc.DefaultData()

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		td.AuthUser = &models.User{
			Name: "user name",
		}

		td.CSRFToken = "csrftoken"

		td.Flash = "flash message"

		td.PageTitle = "Title here"

		td.Add("Mapped", "data in map")

		sc.StaticView.Render(w, r, "index.gohtml", td)

	})
}
