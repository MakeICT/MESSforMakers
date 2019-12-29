package controllers

import (
	"fmt"
	"net/http"

	"github.com/golangcollege/sessions"

	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

//ErrorController implements the handlers required for various client and server errors
type ErrorController struct {
	Controller
	ErrorView views.View
}

//Initialize sets up all the required pieces for the controller
func (c *ErrorController) Initialize(cfg *util.Config, um Users, l *util.Logger, s *sessions.Session) error {
	c.setup(cfg, um, l, s)

	c.ErrorView = views.View{}

	if err := c.ErrorView.LoadTemplates("error"); err != nil {
		return fmt.Errorf("Error loading user templates: %v", err)
	}

	return nil
}

//NotFound handles the case when the resource requested by the client cannot be found
func (c *ErrorController) NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		td, err := c.DefaultData(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		td.Add("Error", "Page not found!")

		td.PageTitle = "404 - Page not found"

		if err := c.ErrorView.Render(w, r, "index.gohtml", td); err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

//NotAllowed handles the case when the client id not allowed to access a resource with a certain HTTP method
func (c *ErrorController) NotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		td, err := c.DefaultData(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		td.Add("Error", "Method not allowed!")

		td.PageTitle = "405 - Method Not Allowed"

		if err := c.ErrorView.Render(w, r, "index.gohtml", td); err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
