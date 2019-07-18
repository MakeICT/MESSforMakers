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

func (sc *StaticController) Initialize(cs *session.CookieStore, db *sqlx.DB, l *util.Logger) error {
	sc.setup(cs, db, l)

	tc, err := loadTemplates([]string{"error", "static"})
	if err != nil {
		return fmt.Errorf("Error loading static templates: %v", err)
	}

	sc.StaticView = views.View{TemplateCache: tc}
	return nil
}

func (sc *StaticController) Root() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sc.StaticView.Render(w, r, "index", &views.TemplateData{})
	})
}
