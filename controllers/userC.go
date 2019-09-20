package controllers

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

type UserController struct {
	Controller
	UserView views.View
}

func (uc *UserController) Initialize(cfg *util.Config, cs *session.CookieStore, db *sqlx.DB, l *util.Logger) error {
	uc.setup(cfg, cs, db, l)
	//TODO initialize templates
	//sync.RunOnce()?

	uc.UserView = views.View{}

	if err := uc.UserView.LoadTemplates([]string{"error", "user"}); err != nil {
		return fmt.Errorf("Error loading user templates: %v", err)
	}

	return nil
}
