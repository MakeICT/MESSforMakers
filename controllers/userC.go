package controllers

import (
	"fmt"
	"net/http"
	"strconv"
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
		td, err := uc.DefaultData()
		if err != nil {
			http.Error(w, "could not generate default data", http.StatusInternalServerError)
			return
		}
		td.Add("Form", util.NewForm(nil))
		uc.UserView.Render(w, r, "signup.gohtml", td)
		return
	})
}

// ListUsers generates alist of all users to display. For admin purposes.
func (uc *UserController) ListUsers() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		form := util.NewForm(r.URL.Query())

		page := 1
		p := form.Get("page")
		if p != "" {
			n, err := strconv.ParseFloat(p, 64)
			if err == nil && n != 0 {
				page = int(n)
			}
		}

		sort := form.Get("sort")
		if sort != "" {
			form.PermittedValues(sort, "name", "dob")
			if !form.Valid() {
				sort = ""
			}
		}

		users, err := uc.Users.GetAll(20, page, sort, "asc")
		if err != nil {
			uc.serverError(w, err)
			return
		}

		td, err := uc.DefaultData()
		if err != nil {
			http.Error(w, "could not generate default data", http.StatusInternalServerError)
			return
		}
		td.Add("Users", users)
		uc.UserView.Render(w, r, "users.gohtml", td)
		return
	})
}

//NewUser saves a new user to the database
func (uc *UserController) NewUser() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		form := util.NewForm(r.PostForm)

		form.Required("name", "email", "email2", "password", "password2", "dob.mm", "dob.dd", "dob.yyyy", "phone")
		form.RequiredIf("membershipoption", r.FormValue("membersignup") == "on")
		form.PermittedValues("membershipoption", "1", "2", "3", "4", "5", "6")
		form.MatchField("email", "email2")
		form.MatchField("password", "password2")
		form.MinLength("password", 4)
		form.MaxLength("name", 255)
		form.MaxLength("email", 255)
		form.MaxLength("phone", 15)
		form.MatchPattern("email", util.EmailRegEx)
		form.MatchPattern("phone", util.PhoneRegEx)

		dob, err := time.Parse("01-02-2006", fmt.Sprintf("%s-%s-%s", r.FormValue("dob.mm"), r.FormValue("dob.dd"), r.FormValue("dob.yyyy")))
		if err != nil {
			form.Errors.Add("dob", "Could not recognize date")
		}

		var ms, mo int
		ms = 1
		if r.FormValue("membersignup") == "on" {
			//The error from Atoi is ignored because the value has already been confirmed to be the string 1, 2, 3, 4, 5, or 6
			mo, _ = strconv.Atoi(r.FormValue("membershipoption"))
			ms = 1
		}

		if !form.Valid() {
			//TODO add flash message that there were errors
			td, err := uc.DefaultData()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			//To minimize the number of times the plaintext password is passed back and forth, remove it before responding
			form.Set("password", "")
			form.Set("password2", "")

			td.Add("Form", form)
			uc.UserView.Render(w, r, "signup.gohtml", td)
			return
		}

		//TODO add gorilla/schema to scan a form directly into a struct

		u := &models.User{
			Name:             r.FormValue("name"),
			Email:            r.FormValue("email"),
			Password:         r.FormValue("password"),
			DOB:              dob,
			Phone:            r.FormValue("phone"),
			TextOK:           r.FormValue("oktotext") == "on",
			MembershipStatus: ms,
			MembershipOption: mo,
		}

		err = uc.Users.Create(u)
		if err != nil {

			td, err2 := uc.DefaultData()
			if err2 != nil {
				http.Error(w, "could not generate default data", http.StatusInternalServerError)
				return
			}

			form.Set("password", "")
			form.Set("password2", "")

			form.Errors.Add("saveError", fmt.Sprintf("Could not save user: %s", err.Error()))

			td.Add("Form", form)
			uc.UserView.Render(w, r, "signup.gohtml", td)
			return
		}

		redir := fmt.Sprintf("http://%s:%d/user/%d", uc.AppConfig.App.Host, uc.AppConfig.App.Port, u.ID)
		uc.Logger.Debugf("Redirecting to: %s", redir)

		http.Redirect(w, r, redir, http.StatusSeeOther)
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
