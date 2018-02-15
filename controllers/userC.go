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

type Address struct {
	Line1 string
	Line2 string
	City  string
	State string
	Zip   string
}

type Ice struct {
	Name     string `schema:"name"`
	Phone    string `schema:"phone"`
	Relation string `schema:"relation"`
}

type Person struct {
	FirstName string `schema:"firstname"`
	LastName  string `schema:"lastname"`
	Address   Address
	Phone     string `schema:"phone"`
	OfAge     bool   `schema:"ofage"`
	Guardian  string `schema:"guardian"`
	Email     string `schema:"email"`
	Password  string `schema:"-"`
	Ice       Ice    `schema:"ice"`
}

type Login struct {
	Email         string `schema:"email"`
	EmailCheck    string `schema:"emailcheck"`
	Password      string `schema:"password"`
	PasswordCheck string `schema:"passwordcheck"`
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

			person := new(Person)
			body := struct{ Person *Person }{Person: person}
			// body := ""
			if err := views.User.New.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		} else if r.Method == "POST" {

			person, _ := parseSignupForm(r, c.Decoder)
			//TODO if email submitted is not valid, or does not match emailcheck, it should not be returned to fill form
			//TODO if password submitted is not valid, or does not match passwordcheck, it should not be returned to fill form
			body := struct{ Person *Person }{Person: person}

			//TODO need a Person.Validate() function to check a populated User object for correctness, that would be a database function? The Person object itself should probably be defined in the model, not here.
			// the Person object should be passed around, so if there are errors, it can be reused to populate the form again
			//TODO how to communicate errors
			if err := views.User.New.Render(w, body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func parseSignupForm(r *http.Request, d *schema.Decoder) (*Person, *Login) {
	if err := r.ParseForm(); err != nil {
		//TODO handle this error
	}

	p := new(Person)
	l := new(Login)
	if err := d.Decode(p, r.PostForm); err != nil {
		//TODO handle this error
	}
	if err := d.Decode(l, r.PostForm); err != nil {
		//TODO handle this error
	}

	return p, l
}
