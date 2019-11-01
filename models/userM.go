package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Address is to nest an address and keep the User struct cleaner, can also be reused elsewhere
type Address struct {
	Line1 string
	Line2 string
	City  string
	State string
	Zip   string
}

// Ice holds the In Case of Emergence information for a user
type Ice struct {
	Name     string `schema:"name"`
	Phone    string `schema:"phone"`
	Relation string `schema:"relation"`
}

// User holds all the details about a single person
type User struct {
	ID               int    `db:"id"`
	FirstName        string `schema:"firstname" db:"first_name"`
	LastName         string `schema:"lastname" db:"last_name"`
	Address          Address
	Phone            string     `schema:"phone" db:"phone"`
	OfAge            bool       `schema:"ofage"`
	DOB              *time.Time `schema:"dob"`
	Guardian         string     `schema:"guardian"`
	Ice              Ice        `schema:"ice"`
	Email            string     `schema:"email" db:"username"`
	EmailCheck       string     `schema:"emailcheck"`
	Password         string     `schema:"password" db:"password"`
	PasswordCheck    string     `schema:"passwordcheck"`
	Authorized       bool
	TextOK           bool `db:"text_ok"`
	MembershipStatus int  `db:"membership_status_id"`
	MembershipOption int  `db:"membership_option"`
	RBACRole         int  `db:"rbac_role_id"`
}

//ErrNotAuthorized is used when the user is not authorized to perform the requested action
// TODO Fill out this list. review all error handling in app and controller
var ErrNotAuthorized = errors.New("authorization failed")

//ErrNoRecord is returned when there are no records to return, but a record is required.
var ErrNoRecord = errors.New("no matching records")

// UserModel stores the database handle and any other globals needed for the database methods
// All user related DB methods will be defined on this model
type UserModel struct {
	DB *sqlx.DB
}

//Get one user (need user ID populated)
func (um *UserModel) Get(id int) (*User, error) {
	if !(id > 0) {
		return nil, fmt.Errorf("Did not recognize user id")
	}
	q := um.DB.Rebind("SELECT id, name, username, dob, phone FROM member WHERE id = ?")
	user := &User{}
	err := um.DB.Get(user, q, id)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve user: %v", err)
	}
	return user, nil
}

//GetAll returns "count" many users, starting "offset" users from the beginning
func (um *UserModel) GetAll(count, page int, sortBy, direction string) ([]*User, error) {
	offset := (page - 1) * count //for page 1 the offset should be 0, etc.
	q := um.DB.Rebind(`
		SELECT 
			id, 
			name, 
			username, 
			dob, 
			phone 
		FROM 
			member 
		ORDER BY 
			name 
		LIMIT 
			? 
		OFFSET 
		?
	`)
	rows, err := um.DB.Queryx(q, count, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var u User
		if err := rows.StructScan(&u); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func generateKey(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// TODO add some error checking to ensure that the OS CSPRNG has not failed in any way

	return base64.URLEncoding.EncodeToString(b), err
}

// OriginateSession takes a user ID and starts a session by creating an authkey and storing that in the session table,
// then returning that key to be used in the session cookie
func (um *UserModel) OriginateSession(id int) (string, error) {
	//generate crypto random key
	key, err := generateKey(32)
	if err != nil {
		return "", err
	}

	//store key in database
	query := "INSERT INTO session (userid, authtoken, loginDate, lastSeenDate) VALUES ($1, $2, $3, $4)"

	//TODO check that the user exists first?
	// TODO make the fields datetime not date
	_, err = um.DB.Exec(query, id, key, "1-1-1970", "1-1-1970")

	if err != nil {
		return "", err
	}

	//return key
	return key, nil
}

// SessionLookup searches for a session in the database, and makes sure that it's not deleted or expired.
//TODO
func (um *UserModel) SessionLookup(id int, auth string) (*User, error) {

	return nil, nil

}

//Create user (need user details populated)
func (um *UserModel) Create(u *User) error {

	var q string
	var guestStatus int
	var guestRole int

	//TODO sanitize strings to prevent SQL injection

	//Fetch the ID for guest status from DB
	q = "SELECT id FROM membership_status WHERE name = 'guest'"
	um.DB.Get(&guestStatus, q)

	//fetch the ID for the most restricted role from DB
	q = "SELECT id FROM rbac_role WHERE name = 'guest'"
	um.DB.Get(&guestRole, q)

	//TODO calculate membership_expires
	q = um.DB.Rebind(`
	INSERT INTO member 
		(name, username, password, dob, phone, membership_status_id, rbac_role_id, created_at, updated_at)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id`)
	var id int
	fmt.Printf("%+v\n", u)
	err := um.DB.Get(
		&id,
		q,
		u.FirstName,
		u.Email,
		u.Password,
		u.DOB,
		u.Phone,
		u.MembershipStatus,
		1,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

//Update user (need user details populated)
func (um *UserModel) Update(*User) error {
	return nil
}

//Delete user (need user ID populated)
func (um *UserModel) Delete(*User) error {
	return nil
}

//CheckPassword returns true only if the password matches the stored password for the user
func (um *UserModel) CheckPassword(username, password string) (int, error) {
	return 1, nil
}

//define database helper functions here

//define Valuers and Scanners for any user related custom types here.
