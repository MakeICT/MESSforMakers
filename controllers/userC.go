/*
 MESS for Makers - An open source member and event management platform
    Copyright (C) 2017  Sam Schurter

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package controllers

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/views"
)

//separate files for each individual controller
//a controller is a struct, with methods defined on the struct for each action

// the struct defines what the controller needs to be able to pass into any given page it needs to render
type UserController struct {
	Controller
	DB *sqlx.DB
}

// setup function that stores the database pool in the controller, or other things if necessary
func User(db *sqlx.DB) UserController {
	return UserController{DB: db}
}

// a handler
func (c *UserController) Index() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//needs to be a slice of User from the database
		//pagination needs to be built in from the start. get from query param if
		//there and store in cookie. else, get from cookie
		users, _ := models.GetAllUsers(c.DB, 10, 0)

		if err := views.User.Index.Render(w, users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (c *UserController) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			body := ""
			if err := views.User.New.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		} else if r.Method == "POST" {
			body := "Arrived via post request. Not set up to save in DB yet"
			// need a User.Validate() function to check a populated User object for correctness, that would be a database function.
			// the User object should be passed around, so if there are errors, it can be reused to populate the form again
			if err := views.ErrorPage.Index.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
