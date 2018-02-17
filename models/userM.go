package models

import (
	"net/http"

	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
)

//TODO abstract into the main model file, this could be used by other controllers
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

type User struct {
	ID            int
	FirstName     string `schema:"firstname"`
	LastName      string `schema:"lastname"`
	Address       Address
	Phone         string `schema:"phone"`
	OfAge         bool   `schema:"ofage"`
	Guardian      string `schema:"guardian"`
	Ice           Ice    `schema:"ice"`
	Email         string `schema:"email"`
	EmailCheck    string `schema:"emailcheck"`
	Password      string `schema:"password"`
	PasswordCheck string `schema:"passwordcheck"`
}

//only needs strings, and no nested structs.  One message string can be used to describe any address or ICE issues
type UserErrors struct {
	FirstName string
	LastName  string
	Address   string
	Phone     string
	OfAge     string
	Guardian  string
	Ice       string
	Email     string
	Password  string
}

//get one user (need user ID populated)
func (u *User) getUser(db *sqlx.DB) error {
	return nil
}

//get "count" many users, starting "offset" users from the beginning
func GetAllUsers(db *sqlx.DB, count, offset int) ([]User, error) {
	if err := db.Ping(); err != nil {
		return nil, err
	}
	rows, err := db.Queryx("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.StructScan(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// handles getting the form data out of the http.Request and into an easier format to work with.
func (u *User) ParseSignupForm(r *http.Request, d *schema.Decoder) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	if err := d.Decode(u, r.PostForm); err != nil {
		return err
	}
	return nil
}

// Validate a user object, return a struct of strings with error messages if it fails.
// Destructive function, modifies the user if there are irreconcilable errors with password or email.
func (u *User) ValidateUser() *UserErrors {
	ue := new(UserErrors)
	errorsFound := false

	//for simple non-regex validators, implement directly for now.
	if len(u.FirstName) < 2 {
		ue.FirstName = "Please enter your full name"
		errorsFound = true
	}

	if len(u.LastName) < 2 {
		ue.LastName = "Please enter your full name"
		errorsFound = true
	}

	// extract regex validators into helper functions, then store the errors into the error struct.
	if err := validAddress(u.Address); err != nil {
		ue.Address = err.Error()
		errorsFound = true
	}

	if err := validPhone(u.Phone); err != nil {
		ue.Address = err.Error()
		errorsFound = true
	}

	if u.OfAge == false && u.Guardian == "" {
		ue.Guardian = "If you are under 18, you must have the permission of a parent or legal guardian."
		errorsFound = true
	}

	if err := validIce(u.Ice); err != nil {
		ue.Ice = err.Error()
		errorsFound = true
	}

	if err := validEmail(u.Email, u.EmailCheck); err != nil {
		//if the emails are invalid for any reason, then delete any value in EmailCheck so it does not render back in the form for correction
		u.EmailCheck = ""
		ue.Email = err.Error()
		errorsFound = true
	}

	if err := validPassword(u.Password, u.PasswordCheck); err != nil {
		//if the passwords are invalid for any reason, clear both so the operator can re-enter
		u.Password = ""
		u.PasswordCheck = ""
		ue.Password = err.Error()
		errorsFound = true
	}

	if errorsFound {
		return ue
	}

	return nil
}

//create user (need user details populated)
//TODO actually store the user in the database. Nil return implies success.
func (u *User) CreateUser(db *sqlx.DB) error {
	return nil
}

//update user (need user details populated)
func (u *User) updateUser(db *sqlx.DB) error {
	return nil
}

//delete user (need user ID populated)
func (u *User) deleteUser(db *sqlx.DB) error {
	return nil
}

//define database helper functions here
// TODO implement all validators
func validAddress(a Address) error {
	//import errors; return errors.New(string)
	//then return can access err via err.Error()
	return nil
}

func validPhone(p string) error {
	return nil
}

func validIce(i Ice) error {
	return nil
}

func validEmail(e string, ec string) error {
	return nil
}

func validPassword(p string, pw string) error {
	return nil
}

//define Valuers and Scanners for any user related custom types here.
