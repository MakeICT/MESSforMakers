// Package controllers provides handlers that can be mapped to routes and defines interfaces that define
// what is required by the handlers.
package controllers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/golangcollege/sessions"

	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

// Users interface defines the methods that a Users model must fulfill. Allows mocking with a fake database for testing.
type Users interface {
	Get(int) (*models.User, error)
	GetAll(int, int, string, string) ([]*models.User, error)
	Create(*models.User) error
	Update(*models.User) error
	Delete(*models.User) error
	Login(string, string, string, string) (int, string, error)
	SessionLookup(int, string) (*models.User, error)
}

// Controller is a struct Struct to store pointer to cookiestore, database, and logger and any other things common to many controllers
type Controller struct {
	Users     Users
	Logger    *util.Logger
	AppConfig *util.Config
	Session   *sessions.Session
}

// method to create a new struct and store the information from the app, passed as args
// Requiring the information passed as args avoids imports loop
// the general controller constructor will be called by the specific controller constructors
// the specific controller constructors can then initialize and embed their own templates.
func (c *Controller) setup(cfg *util.Config, um Users, l *util.Logger, s *sessions.Session) {
	c.Users = um
	c.Logger = l
	c.AppConfig = cfg
	c.Session = s
}

// DefaultData ia the method to generate required default template data and return template object
func (c *Controller) DefaultData(r *http.Request) (*views.TemplateData, error) {
	td := &views.TemplateData{}
	td.Root = fmt.Sprintf("http://%s:%d/", c.AppConfig.App.Host, c.AppConfig.App.Port)
	td.Flash = c.Session.PopString(r, "flash")
	return td, nil
}

func (c *Controller) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	c.Logger.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (c *Controller) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (c *Controller) notFound(w http.ResponseWriter) {
	c.clientError(w, http.StatusNotFound)
}
