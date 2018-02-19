package models

import (
	"fmt"
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
	FirstName     string `schema:"firstname" db:"first_name"`
	LastName      string `schema:"lastname" db:"last_name"`
	Address       Address
	Phone         string `schema:"phone" db:"phone"`
	OfAge         bool   `schema:"ofage"`
	DOB           string `schema:"dob"`
	Guardian      string `schema:"guardian"`
	Ice           Ice    `schema:"ice"`
	Email         string `schema:"email" db:"username"`
	EmailCheck    string `schema:"emailcheck"`
	Password      string `schema:"password" db:"password"`
	PasswordCheck string `schema:"passwordcheck"`
}

//only needs strings, and no nested structs.  One message string can be used to describe any address or ICE issues
type UserErrors struct {
	FirstName string
	LastName  string
	Address   string
	Phone     string
	OfAge     string
	DOB       string
	Guardian  string
	Ice       string
	Email     string
	Password  string
}

//get one user (need user ID populated)
func (u *User) GetUser(db *sqlx.DB) error {

	query := "SELECT id, first_name, last_name, username, phone, password FROM member WHERE id = $1"

	if err := db.Get(u, query, u.ID); err != nil {
		return err
	}

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

	if err := validDob(u.DOB); err != nil {
		// ue.DOB = err.Error()
		// errorsFound = true
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

	var query string
	var guestStatus int
	var guestRole int

	//TODO sanitize strings to prevent SQL injection

	//Fetch the ID for guest status from DB
	query = "SELECT id FROM membership_status WHERE name = 'guest'"
	db.Get(&guestStatus, query)

	//fetch the ID for the most restricted role from DB
	query = "SELECT id FROM rbac_role WHERE name = 'guest'"
	db.Get(&guestRole, query)

	query = `INSERT INTO member (
		first_name, 
		last_name,
		username, 
		password, 
		dob, 
		phone, 
		membership_status_id, 
		rbac_role_id 
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	//replace queryrow/scan with sqlx get. needs proper struct tags to work
	err := db.QueryRow(query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
		"1-1-1970",
		u.Phone,
		guestStatus,
		guestRole,
	).Scan(&u.ID)
	if err != nil {
		return err
	}
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

func validDob(d string) error {
	if d == "" {
		return fmt.Errorf("Must enter a valid date of birth")
	}
	return nil
}

func validEmail(e string, ec string) error {
	if e == "" {
		return fmt.Errorf("Must enter a valid email address")
	}
	if e != ec {
		return fmt.Errorf("Email addresses must match")
	}
	return nil
}

func validPassword(p string, pc string) error {
	if p == "" {
		return fmt.Errorf("Must enter a valid password")
	}
	if p != pc {
		return fmt.Errorf("Passwords must match")
	}
	return nil
}

//define Valuers and Scanners for any user related custom types here.
