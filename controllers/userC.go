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
	// "fmt"
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

//Body struct to contain all the various fields that may be needed to render the templates.
type Body struct {
	User          *models.User
	ErrorMessages *models.UserErrors
	Login         *models.Login
}

// setup function that stores the database pool in the controller, or other things if necessary
// returns a struct that holds the user-related data. Can be injected for testing.
func User(db *sqlx.DB, cs *session.CookieStore) UserController {
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

			// always need a body, even if it's empty. The template can check for empty fields, but not missing structs
			body := Body{}

			if err := views.User.New.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if r.Method == "POST" {

			//A user object can be passed around to contain and preserve partial form data that needs revision and error correction.
			user := &models.User{}

			//the model should not be handling things like errors and http. Send the model only the minimum necessary for parsing the form.
			// TODO Refactor to get the http.Request out of the model
			if err := user.ParseSignupForm(r, c.Decoder); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//when the form has parsed, begin validating the data, and filling the form with whatever can be salvaged from an invalid submission.
			// TODO the model should be returning better error messages, so that better decisions can be made based on WHY the validation failed.
			if errMessages := user.ValidateUser(); errMessages != nil {
				body := Body{
					User:          user,
					ErrorMessages: errMessages,
				}
				if err := views.User.New.Render(w, body); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			//only create the user in the database after all logic-level validation has passed. It's still possible to fail at this point
			// because the database will enforce types.
			//TODO since it can still fail at this point due to DB enforced rules, replace internal server error with a more friendly rendering of the signup form.
			//Nil return only implies success, the existence of the user.ID is more important
			if err := user.CreateUser(c.DB); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//if the user is created successfully, there is nothing to render here, and it should show the user their new profile page.
			//use StatusSeeOther for redirect.
			http.Redirect(w, r, "/user/"+strconv.Itoa(user.ID), http.StatusSeeOther)
		}
	}
}

//Show a single user.
func (c *UserController) Show() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// the router will only call this function if there is a numeric ID number in the URL,
		// but it will need to be converted from a string to an int.
		// TODO if it can't convert from string to int, it should do something nicer then 500
		id_string := mux.Vars(r)["id"]
		id, err := strconv.Atoi(id_string)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//the model needs a user with a minimum of an ID to fetch the data to render the user's page
		// TODO the model should return better errors, if the user is not found, the error should be 404 not 500
		u := &models.User{ID: id}
		if err := u.GetUser(c.DB); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body := Body{User: u}

		if err := views.User.Show.Render(w, body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func (c *UserController) Edit() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			//TODO this code is repeated exactly from the Show() method and twice here, and will be in any method that requires the user's ID from the URL. DRY it.
			id_string := mux.Vars(r)["id"]
			id, err := strconv.Atoi(id_string)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// In order to pre-fill the form, the user has to be retrieved.
			u := &models.User{ID: id}
			if err := u.GetUser(c.DB); err != nil {
				// TODO again, if the user can't be retrieved, handle the error better and more politely.
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			body := Body{User: u}

			if err := views.User.Edit.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if r.Method == "POST" {

			//A user object can be passed around to contain and preserve partial form data that needs revision and error correction.
			id_string := mux.Vars(r)["id"]
			id, err := strconv.Atoi(id_string)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//the model should not be handling things like errors and http. Send the model only the minimum necessary for parsing the form.
			user := &models.User{ID: id}
			if err := user.ParseSignupForm(r, c.Decoder); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//when the form has parsed, begin validating the data, and filling the form with whatever can be salvaged from an invalid submission.
			if errMessages := user.ValidateUser(); errMessages != nil {
				body := Body{
					User:          user,
					ErrorMessages: errMessages,
				}
				if err := views.User.Edit.Render(w, body); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
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

//Displays the form allowing the user to log in, and handles the form data and login procedure when the form is submitted.
func (c *UserController) Login() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var err error

		if r.Method == "GET" {

			body := Body{}

			if err = views.User.Login.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if r.Method == "POST" {

			login := &models.Login{}

			//get post data into login object
			if err = login.ParseLoginForm(r, c.Decoder); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//validate password "hash" for username
			user := &models.User{}
			if user, err = login.ValidateLogin(); err != nil {

				//if the user comes back invalid for any reason, use the username from the submitted form, but nothing else
				body := Body{
					Login: &models.Login{
						Username: login.Username,
					},
					ErrorMessages: &models.UserErrors{
						Login: "Invalid username or password",
					},
				}

				if err = views.User.Login.Render(w, body); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			//If everything comes back valid, set the cookie and start a session
			if err = setLoginCookie(w, r, c, user); err != nil {
				//TODO failed to log in. Should indicate that somehow.
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			//use StatusSeeOther for redirect.
			// TODO this should ideally take the user back where they were before logging in, NOT directly back to the userpage
			http.Redirect(w, r, "/user/"+strconv.Itoa(user.ID), http.StatusSeeOther)
		}
	}
}

// this should set start a session in the database by generating an auth token and storing it, then setting a cookie with that auth token and the user's ID
func setLoginCookie(w http.ResponseWriter, r *http.Request, c *UserController, u *models.User) error {

	// originate session
	// method needs to generate an auth token, store the token in the DB with the user ID and return the token and an error
	authToken, err := u.OriginateSession(c.DB)
	if err != nil {
		return err
	}
	//if there is no error,
	session, err := c.CookieStore.Store.Get(r, "mess-data")
	if err != nil {
		return err
	}

	session.Values["userid"] = u.ID
	session.Values["authtoken"] = authToken

	if err = session.Save(r, w); err != nil {
		return err
	}

	return nil
}
