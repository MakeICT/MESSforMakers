package controllers

import (
	"github.com/jmoiron/sqlx"

	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

type UserController struct {
	Controller
	UserView views.View
}

func (uc *UserController) Initialize(cs *session.CookieStore, db *sqlx.DB, l *util.Logger) {
	uc.setup(cs, db, l)
	//TODO initialize templates
	uc.UserView = views.View{}
}
