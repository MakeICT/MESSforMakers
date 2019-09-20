package controllers

import (
	"fmt"
	"net/http"

	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

type StaticController struct {
	Controller
	StaticView views.View
}

func (sc *StaticController) Initialize(cfg *util.Config, cs *session.CookieStore, um Users, l *util.Logger) error {
	sc.setup(cfg, cs, um, l)
	sc.StaticView = views.View{}

	if err := sc.StaticView.LoadTemplates([]string{"error", "static"}); err != nil {
		return fmt.Errorf("Error loading static templates: %v", err)
	}

	return nil
}

func (sc *StaticController) Root() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		td := &views.TemplateData{}

		if err := sc.AddDefaultData(td); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sc.StaticView.Render(w, r, "index", td)

	})
}
