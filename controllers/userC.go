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
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"

	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/views"
)

//separate files for each individual controller
//a controller is a struct, with methods defined on the struct for each action

// the struct defines what the controller needs to be able to pass into any given page it needs to render
type UserController struct {
	Controller
	DB      *sqlx.DB
	Decoder *schema.Decoder
}

// setup function that stores the database pool in the controller, or other things if necessary
func User(db *sqlx.DB) UserController {
	//TODO make a map for holding state names and abbreviations, prepopulate it at setup time with a list read from file.
	// list available here: https://gist.github.com/mshafrir/2646763
	return UserController{DB: db, Decoder: schema.NewDecoder()}
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

			user := new(models.User)
			errMessages := new(models.UserErrors)
			body := struct {
				ErrorMessages *models.UserErrors
				Person        *models.User
			}{

				ErrorMessages: errMessages,
				Person:        user,
			}
			if err := views.User.New.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		} else if r.Method == "POST" {

			user := new(models.User)
			if err := user.ParseSignupForm(r, c.Decoder); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if errMessages := user.ValidateUser(); errMessages != nil {
				fmt.Printf("%+v\n", errMessages)
				body := struct {
					Person        *models.User
					ErrorMessages *models.UserErrors
				}{
					Person:        user,
					ErrorMessages: errMessages,
				}
				if err := views.User.New.Render(w, body); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			if err := user.CreateUser(c.DB); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/user/"+strconv.Itoa(user.ID), http.StatusSeeOther)
		}
	}
}
