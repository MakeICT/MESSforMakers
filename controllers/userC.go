package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/golangcollege/sessions"
	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/util"
	"github.com/makeict/MESSforMakers/views"
)

//UserController implements the handlers required for user management
type UserController struct {
	Controller
	UserView views.View
}

//Initialize performs the required setup for a user controller
func (uc *UserController) Initialize(cfg *util.Config, um Users, l *util.Logger, s *sessions.Session) error {
	uc.setup(cfg, um, l, s)

	uc.UserView = views.View{}

	if err := uc.UserView.LoadTemplates("user"); err != nil {
		return fmt.Errorf("Error loading user templates: %v", err)
	}

	return nil
}

//SignupForm displays the signup form
func (uc *UserController) SignupForm() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td, err := uc.DefaultData(r)
		if err != nil {
			http.Error(w, "could not generate default data", http.StatusInternalServerError)
			return
		}
		td.Add("Form", util.NewForm(nil))

		if err := uc.UserView.Render(w, r, "signup.gohtml", td); err != nil {
			uc.serverError(w, err)
			return
		}
	})
}

// List generates alist of all users to display. For admin purposes.
func (uc *UserController) List() func(http.ResponseWriter, *http.Request) {
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

		td, err := uc.DefaultData(r)
		if err != nil {
			http.Error(w, "could not generate default data", http.StatusInternalServerError)
			return
		}
		td.Add("Users", users)
		if err := uc.UserView.Render(w, r, "users.gohtml", td); err != nil {
			uc.serverError(w, err)
			return
		}
	})
}

//Show gets the parameter from the url and gets the details for that user from the database
func (uc *UserController) Show() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := util.IntOK(vars["id"], 1, math.MaxInt8)
		if !ok {
			uc.clientError(w, http.StatusBadRequest)
			return
		}
		user, err := uc.Users.Get(id)
		if err != nil {
			uc.serverError(w, err)
			return
		}

		td, err := uc.DefaultData(r)
		if err != nil {
			uc.serverError(w, err)
		}

		td.Add("User", user)
		uc.Logger.Printf("user: %+v", user)

		if err := uc.UserView.Render(w, r, "show.gohtml", td); err != nil {
			uc.serverError(w, err)
			return
		}
	})
}

//New saves a new user to the database
func (uc *UserController) New() func(http.ResponseWriter, *http.Request) {
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

		//TODO should recognize non-zero-padded months and days e.g. 6 for june
		//would help to confirm int on all three fields, and supply reasonable ranges (1-12, 1-31, and 1900-curyear)
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

		/*
			if u.OfAge == false && u.Guardian == "" {
				ue.Guardian = "If you are under 18, you must have the permission of a parent or legal guardian."
				errorsFound = true
			}
		*/

		if !form.Valid() {
			td, err := uc.DefaultData(r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			//TODO this erases any existing flash messages. refactor to have flash be a slice of strings
			td.Flash = "There were errors saving the form"

			//To minimize the number of times the plaintext password is passed back and forth, remove it before responding
			form.Set("password", "")
			form.Set("password2", "")

			td.Add("Form", form)
			uc.UserView.Render(w, r, "signup.gohtml", td)
			return
		}

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

			td, err2 := uc.DefaultData(r)
			if err2 != nil {
				http.Error(w, "could not generate default data", http.StatusInternalServerError)
				return
			}

			form.Set("password", "")
			form.Set("password2", "")

			td.Flash = fmt.Sprintf("Could not save user: %s", err.Error())

			td.Add("Form", form)
			uc.UserView.Render(w, r, "signup.gohtml", td)
			return
		}

		uc.Session.Put(r, "flash", "Successfully saved user!")

		http.Redirect(w, r, fmt.Sprintf("http://%s:%d/user/%d", uc.AppConfig.App.Host, uc.AppConfig.App.Port, u.ID), http.StatusSeeOther)
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

/**********************************************************/ //
/*
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
*/
