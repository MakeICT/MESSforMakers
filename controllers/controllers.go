// Package controllers provides handlers that can be mapped to routes and defines interfaces that define
// what is required by the handlers.
package controllers

import (
	"net/http"

	"github.com/makeict/MESSforMakers/views"
)

// Controller is a generic type that allows methods to be defined across all controllers, since Controller is embedded in all controllers
type Controller struct {
	AddDefaultData func(*views.TemplateData) *views.TemplateData
}

// NotImplementedController returns an empty controller struct, allowing a route builder to define routes before there is a corresponding handler.
func NotImplementedController() Controller {
	return Controller{}
}

// None returns a message indicating that the route does not exist yet, allowing routing tables to be built without needing all the handler code in place.
func (c *Controller) None(route string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body := "This route has not been implemented yet: " + route

		if err := views.ErrorPage.Index.Render(w, body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
