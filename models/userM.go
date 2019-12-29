package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/makeict/MESSforMakers/util"
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
	ID               int        `db:"id"`
	FirstName        string     `schema:"firstname" db:"first_name"`
	LastName         string     `schema:"lastname" db:"last_name"`
	Address          Address    `schema:"-"`
	Phone            string     `schema:"phone" db:"phone"`
	OfAge            bool       `schema:"ofage"`
	DOB              *time.Time `schema:"dob"`
	Guardian         string     `schema:"guardian"`
	Ice              Ice        `schema:"ice"`
	Email            string     `schema:"email" db:"username"`
	EmailCheck       string     `schema:"emailcheck"`
	Password         string     `schema:"password" db:"password"`
	PasswordCheck    string     `schema:"passwordcheck"`
	Authorized       bool       `schema:"-"`
	TextOK           bool       `db:"text_ok"`
	MembershipStatus int        `db:"membership_status_id"`
	MembershipOption int        `db:"membership_option"`
	RBACRole         int        `db:"rbac_role_id"`
}

//UserSession is used for storing the details of a session retrieved from the DB
type UserSession struct {
	ID         int        `db:"member_id"`
	AuthKey    string     `db:"authtoken"`
	Originated *time.Time `db:"originated"`
	LastSeen   *time.Time `db:"last_seen"`
	IP         string     `db:"last_ip"`
	UserAgent  string     `db:"agent"`
}

//ErrNotAuthorized is used when the user is not authorized to perform the requested action
var ErrNotAuthorized = errors.New("authorization failed")

//ErrNoRecord is returned when there are no records to return, but a record is required.
var ErrNoRecord = errors.New("no matching records")

//ErrBadUsernamePassword is used when the user cannot be logged in due to no matching user in the database or because the wrong password was used.
var ErrBadUsernamePassword = errors.New("no matching user or password does not match")

//ErrSessionExpired is returned if the user's session is no longer valid
var ErrSessionExpired = errors.New("user session has expired")

// UserModel stores the database handle and any other globals needed for the database methods
// All user related DB methods will be defined on this model
type UserModel struct {
	DB       *sqlx.DB
	HashCost int
	logger   *util.Logger
}

//ContextKey will be the type used to retrieve values from a context
type ContextKey int

//ContextkeyUser is used to retrieve the user struct from the context if it exists
const ContextkeyUser ContextKey = 1

//Get one user (need user ID populated)
func (um *UserModel) Get(id int) (*User, error) {
	if !(id > 0) {
		return nil, fmt.Errorf("Did not recognize user id")
	}
	q := um.DB.Rebind("SELECT id, first_name, last_name, username, dob, phone FROM member WHERE id = ?")
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
			first_name, 
			last_name, 
			username, 
			dob, 
			phone
		FROM 
			member 
		ORDER BY 
			%s 
		LIMIT 
			? 
		OFFSET 
			?
	`)
	orderBy := ""
	direction = strings.ToUpper(direction)
	switch sortBy {
	case "fname":
		orderBy = "first_name " + direction
	case "lname":
		orderBy = "last_name " + direction
	case "dob":
		orderBy = "dob " + direction
	default:
		orderBy = "id asc"
	}

	q = fmt.Sprintf(q, orderBy)

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

//Create user (need user details populated)
func (um *UserModel) Create(u *User) error {

	var q string
	var guestStatus int
	var guestRole int

	//Fetch the ID for guest status from DB
	q = "SELECT id FROM membership_status WHERE name = 'guest'"
	um.DB.Get(&guestStatus, q)

	//fetch the ID for the most restricted role from DB
	q = "SELECT id FROM rbac_role WHERE name = 'guest'"
	um.DB.Get(&guestRole, q)

	q = um.DB.Rebind(`
	INSERT INTO member 
		(first_name, last_name, username, password, dob, phone, membership_status_id, rbac_role_id, created_at, updated_at)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id`)
	var id int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), um.HashCost)
	if err != nil {
		return err
	}

	err = um.DB.Get(
		&id,
		q,
		u.FirstName,
		u.LastName,
		u.Email,
		hashedPassword,
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

// Login takes a username and password and starts a session by creating an authkey and storing that in the session table,
// then returning that key to be used in the session cookie
func (um *UserModel) Login(u, p, ip, ua string) (int, string, error) {

	query := "SELECT id, password FROM member WHERE username = $1"
	var id int
	var password []byte
	err := um.DB.QueryRowx(query, u, p).Scan(&id, &password)
	if err == sql.ErrNoRows {
		query = "INSERT INTO login_log (username, login_status_id)"
		logErr := um.logAttempt(u, "badusername")
		if logErr != nil {
			um.logger.Debugf("could not log login attempt: %v", logErr)
		}
		return 0, "", ErrBadUsernamePassword
	} else if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword(password, []byte(p))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		logErr := um.logAttempt(u, "badpassword")
		if logErr != nil {
			um.logger.Debugf("could not log login attempt: %v", logErr)
		}
		return 0, "", ErrBadUsernamePassword
	} else if err != nil {
		return 0, "", err
	}

	//generate crypto random key
	key, err := generateKey(32)
	if err != nil {
		return 0, "", err
	}

	//store session key in database
	query = `
		INSERT INTO 
			session (member_id, authtoken, originated, last_seen, last_ip, agent) 
		VALUES 
			($1, $2, CURRENT_TIMESTAMP(0), CURRENT_TIMESTAMP(0), $3, $4)
	`

	_, err = um.DB.Exec(query, id, key, ip, ua)

	if err != nil {
		return 0, "", err
	}

	err = um.logAttempt(u, "success")
	if err != nil {
		um.logger.Debugf("could not log login attempt: %v", err)
	}

	//return user id and auth key
	return id, key, nil
}

// SessionLookup searches for a session in the database, and makes sure that it's not deleted or expired.
//if there
func (um *UserModel) SessionLookup(id int, auth string) (*User, error) {

	q := um.DB.Rebind("SELECT authtoken, originated, last_seen, last_ip, agent FROM session WHERE member_id = ? AND authkey = ?")
	r := UserSession{}
	err := um.DB.QueryRowx(q, id, auth).StructScan(r)
	if err == sql.ErrNoRows {
		return nil, ErrNoRecord
	}
	if err != nil {
		return nil, err
	}

	if time.Since(*r.Originated) > time.Hour*(24*30) || time.Since(*r.LastSeen) > time.Hour*(24*7) {
		return nil, ErrSessionExpired
	}

	u, err := um.Get(id)
	if err != nil {
		return nil, err
	}

	return u, nil

}

// SessionDelete removes a specific session from the sessions database given an auth key and user ID.
func (um *UserModel) SessionDelete(id int, auth string) error {
	q := um.DB.Rebind("DELETE FROM session WHERE userid = ? AND authtoken = ?")
	_, err := um.DB.Exec(q, id, auth)
	if err != nil {
		return err
	}
	return nil
}

//SessionUpdate sets the user session's last seen time to now
func (um *UserModel) SessionUpdate(id int, auth string) error {
	q := um.DB.Rebind("UPDATE session SET last_seen = CURRENT_TIMESTAMP(0) WHERE member_id = ? AND authtoken = ?")
	_, err := um.DB.Exec(q, id, auth)
	if err != nil {
		return err
	}
	return nil
}

func (um *UserModel) logAttempt(username, status string) error {
	queryBadUsername := "INSERT INTO login_log (username, login_status_id, created_at) VALUES (?, (SELECT id FROM login_status WHERE name LIKE '%Username'), CURRENT_TIMESTAMP(0))"
	queryBadPassword := "INSERT INTO login_log (username, login_status_id, created_at) VALUES (?, (SELECT id FROM login_status WHERE name LIKE '%Password'), CURRENT_TIMESTAMP(0))"
	querySuccess := "INSERT INTO login_log (username, login_status_id, created_at) VALUES (?, (SELECT id FROM login_status WHERE name LIKE '%Password'), CURRENT_TIMESTAMP(0))"
	q := ""
	switch status {
	case "badpassword":
		q = um.DB.Rebind(queryBadPassword)
	case "badusername":
		q = um.DB.Rebind(queryBadUsername)
	case "success":
		q = um.DB.Rebind(querySuccess)
	default:
		return fmt.Errorf("not a valid login attempt status")
	}
	_, err := um.DB.Exec(q, username)
	if err != nil {
		return err
	}
	return nil
}

//define database helper functions here

func generateKey(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), err
}

//define Valuers and Scanners for any user related custom types here.
