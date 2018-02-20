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

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"

	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/views"
)

//separate files for each individual controller
//a controller is a struct, with methods defined on the struct for each action

// the struct defines what the controller needs to be able to pass into any given page it needs to render
type UserController struct {
	Controller
	DB          *sqlx.DB
	Decoder     *schema.Decoder
	CookieStore *session.CookieStore
}

// setup function that stores the database pool in the controller, or other things if necessary
func User(db *sqlx.DB, cs *session.CookieStore) UserController {
	//TODO make a map for holding state names and abbreviations, prepopulate it at setup time with a list read from file.
	// list available here: https://gist.github.com/mshafrir/2646763
	return UserController{DB: db, Decoder: schema.NewDecoder(), CookieStore: cs}
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

// Create() handles both the request to create, and the posting of the user details.
// Renders the empty form, prefilled form with errors if the form doesn't validate, and redirects to the user's dashboard if the creation is successful.
func (c *UserController) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			//TODO Can this be factored out?  These structs are only here to make the template happy on initial load.
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
				return
			}

		} else if r.Method == "POST" {

			//A user object can be passed around to contain and preserve partial form data that needs revision and error correction.
			user := new(models.User)

			//the model should not be handling things like errors and http. Send the model only the minimum necessary for parsing the form.
			if err := user.ParseSignupForm(r, c.Decoder); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//when the form has parsed, begin validating the data, and filling the form with whatever can be salvaged from an invalid submission.
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

			//only create the user in the database after all logic-level validation has passed. It's still possible to fail at this point
			// because the database will enforce types.
			//TODO since it can still fail at this point, replace internal server error with a more friendly rendering of the signup form.
			//Nil return only implies success, the existence of the user.ID is more important
			if err := user.CreateUser(c.DB); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//use StatusSeeOther for redirect.
			http.Redirect(w, r, "/user/"+strconv.Itoa(user.ID), http.StatusSeeOther)
		}
	}
}

func (c *UserController) Show() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		id_string := mux.Vars(r)["id"]

		u := new(models.User)
		id, err := strconv.Atoi(id_string)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.ID = id

		if err := u.GetUser(c.DB); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body := struct{ Person *models.User }{Person: u}

		if err := views.User.Show.Render(w, body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (c *UserController) Edit() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			id_string := mux.Vars(r)["id"]

			u := new(models.User)
			id, err := strconv.Atoi(id_string)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			u.ID = id

			if err := u.GetUser(c.DB); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ue := new(models.UserErrors)
			body := struct {
				Person        *models.User
				ErrorMessages *models.UserErrors
			}{
				Person:        u,
				ErrorMessages: ue,
			}

			if err := views.User.Edit.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if r.Method == "POST" {

			id_string := mux.Vars(r)["id"]

			//A user object can be passed around to contain and preserve partial form data that needs revision and error correction.
			user := new(models.User)
			id, err := strconv.Atoi(id_string)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			user.ID = id

			//the model should not be handling things like errors and http. Send the model only the minimum necessary for parsing the form.
			if err := user.ParseSignupForm(r, c.Decoder); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//when the form has parsed, begin validating the data, and filling the form with whatever can be salvaged from an invalid submission.
			if errMessages := user.ValidateUser(); errMessages != nil {
				body := struct {
					Person        *models.User
					ErrorMessages *models.UserErrors
				}{
					Person:        user,
					ErrorMessages: errMessages,
				}
				if err := views.User.Edit.Render(w, body); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			//only update the user in the database after all logic-level validation has passed. It's still possible to fail at this point
			// because the database will enforce types.
			//TODO since it can still fail at this point, replace internal server error with a more friendly rendering of the signup form.
			//Nil return only implies success, the existence of the user.ID is more important
			if err := user.UpdateUser(c.DB); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//use StatusSeeOther for redirect.
			http.Redirect(w, r, "/user/"+strconv.Itoa(user.ID), http.StatusSeeOther)
		}
	}
}

func (c *UserController) Login() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			type Login struct {
				Username string
				Password string
				Remember bool
			}
			body := struct {
				Login Login
			}{
				Login: Login{
					Username: "",
					Password: "",
					Remember: false,
				},
			}

			if err := views.User.Login.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if r.Method == "POST" {

			//use StatusSeeOther for redirect.
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}
