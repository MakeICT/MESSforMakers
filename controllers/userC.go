package controllers

import (
	"fmt"
	"net/http"

	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

//UserController implements the handlers required for user management
type UserController struct {
	Controller
	UserView views.View
}

//Initialize performs the required setup for a user controller
func (uc *UserController) Initialize(cfg *util.Config, cs *session.CookieStore, um Users, l *util.Logger) error {
	uc.setup(cfg, cs, um, l)

	uc.UserView = views.View{}

	if err := uc.UserView.LoadTemplates("user"); err != nil {
		return fmt.Errorf("Error loading user templates: %v", err)
	}

	return nil
}

//SignupForm displays the signup form
func (uc *UserController) SignupForm() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "signup form not implemented yet", http.StatusInternalServerError)
		return
	})
}

//NewUser saves a new user to the database
func (uc *UserController) NewUser() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "new user not implemented yet", http.StatusInternalServerError)
		return
	})
}

//LoginForm displays the log in form a
func (uc *UserController) LoginForm() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "login form not implemented yet", http.StatusInternalServerError)
		return
	})
}

//LoginUser performs a login
func (uc *UserController) LoginUser() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "login user not implemented yet", http.StatusInternalServerError)
		return
	})
}

//Logout logs a user out
func (uc *UserController) Logout() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "logout not implemented yet", http.StatusInternalServerError)
		return
	})
}
