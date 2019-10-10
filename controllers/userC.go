package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/makeict/MESSforMakers/models"

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
<<<<<<< HEAD
		td, err := uc.DefaultData()
		if err != nil {
			http.Error(w, "could not generate default data", http.StatusInternalServerError)
			return
		}
		uc.UserView.Render(w, r, "signup.gohtml", td)
		return
=======

		td, err := uc.DefaultData()
		if err != nil {
			http.Error(w, "Could not generate default data", http.StatusInternalServerError)
			return
		}

		uc.UserView.Render(w, r, "signup.gohtml", td)

>>>>>>> e75276c1e8f8fb845994b2b9487a2de665bf54b2
	})
}

//NewUser saves a new user to the database
func (uc *UserController) NewUser() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		dob, err := time.Parse("MM-DD-YYYY", r.FormValue("dob"))
		if err != nil {
			uc.Logger.Debugf("Could not parse DOB (%s), setting to NOW: %s", r.FormValue("dob"), time.Now())
			dob = time.Now()
		}

		u := &models.User{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			DOB:      dob,
			Phone:    r.FormValue("phone"),
			TextOK:   r.FormValue("oktotext") == "on",
		}

		err = uc.Users.Create(u)
		if err != nil {
			td, err := uc.DefaultData()
			if err != nil {
				http.Error(w, "could not generate default data", http.StatusInternalServerError)
				return
			}
			td.Add("User", u)
			uc.UserView.Render(w, r, "signup.gohtml", td)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("%s:%d", uc.AppConfig.App.Host, uc.AppConfig.App.Port), http.StatusSeeOther)
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
