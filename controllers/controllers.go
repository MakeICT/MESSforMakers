// Package controllers provides handlers that can be mapped to routes and defines interfaces that define
// what is required by the handlers.
package controllers

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/jmoiron/sqlx"

	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

// Struct to store pointer to cookiestore, database, and logger

type Controller struct {
	CookieStore *session.CookieStore
	DB          *sqlx.DB
	Logger      *util.Logger
}

// method to create a new struct and store the information from the app, passed as args
// Requiring the information passed as args avoids imports loop
// the general controller constructor will be called by the specific controller constructors
// the specific controller constructors can then initialize and embed their own templates.
func (c *Controller) setup(cs *session.CookieStore, db *sqlx.DB, l *util.Logger) {
	c.CookieStore = cs
	c.DB = db
	c.Logger = l
}

// method to generate required default template data and return template object
func (c *Controller) AddDefaultData(td *views.TemplateData) *views.TemplateData {
	return &views.TemplateData{}
}

func loadTemplates(ff []string) (*template.Template, error) {
	var files []string
	for _, f := range ff {
		fg, err := filepath.Glob(fmt.Sprintf("templates/%s/*.gohtml", f))
		if err != nil {
			return nil, fmt.Errorf("could not find template files: %v", err)
		}
		files = append(files, fg...)
	}
	//fmt.Println(files)
	return template.ParseFiles(files...)
}
