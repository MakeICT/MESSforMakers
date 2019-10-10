package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// UserModel stores the database handle and any other globals needed for the database methods
// All user related DB methods will be defined on this model
type UserModel struct {
	DB *sqlx.DB
}

// User store all information about a user
type User struct {
	ID               int
	Name             string
	Email            string
	Password         string
	DOB              time.Time
	Phone            string
	TextOK           bool
	MembershipStatus int
	MembershipOption int
	RBACRole         int
}

//Get one user (need user ID populated)
func (um *UserModel) Get(int) (*User, error) {
	return nil, nil
}

//GetAll returns "count" many users, starting "offset" users from the beginning
func (um *UserModel) GetAll(count, offset int) ([]User, error) {
	if err := um.DB.Ping(); err != nil {
		return nil, err
	}
	rows, err := um.DB.Queryx("SELECT * FROM users")
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

//Create user (need user details populated)
func (um *UserModel) Create(*User) error {
	return fmt.Errorf("not set up to save user")
}

//Update user (need user details populated)
func (um *UserModel) Update(*User) error {
	return nil
}

//Delete user (need user ID populated)
func (um *UserModel) Delete(*User) error {
	return nil
}

//define database helper functions here

//define Valuers and Scanners for any user related custom types here.
