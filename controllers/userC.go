package controllers

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/views"
)

// UserController embeds the Controller type and stores the data required by User handler
type UserController struct {
	Controller
	DB *sqlx.DB
}

// UserApp defines an interface requiring the Application to provide the necessary data the User handlers will need
type UserApp interface {
	AddDefaultData() func(*views.TemplateData) *views.TemplateData
	DB() *sqlx.DB
}

// User requires an app struct providing data that the User handlers will need, and stores that data in a UserController, which is returned
func User(u UserApp) UserController {
	return UserController{
		DB:             u.DB(),
		AddDefaultData: u.AddDefaultData(),
	}
}

// Index renders a list of all registered users.
func (c *UserController) Index() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//needs to be a slice of User from the database
		//TODO needs pagination. get from query param if
		//there and store in cookie. else, get from cookie
		users, _ := models.GetAllUsers(c.DB, 10, 0)

		//now before rendering, the controller can create a struct with all the data the view needs embedded
		//Then call the function stored in the controller to get the default app data

		if err := views.User.Index.Render(w, users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
