package controllers

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

type StaticController struct {
	Controller
	StaticView views.View
}

func (sc *StaticController) Initialize(cfg *util.Config, cs *session.CookieStore, db *sqlx.DB, l *util.Logger) error {
	sc.setup(cfg, cs, db, l)
	sc.StaticView = views.View{}

	if err := sc.StaticView.LoadTemplates([]string{"error", "static"}); err != nil {
		return fmt.Errorf("Error loading static templates: %v", err)
	}

	return nil
}

func (sc *StaticController) Root() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td := sc.AddDefaultData(&views.TemplateData{})
		sc.StaticView.Render(w, r, "index", td)
	})
}
