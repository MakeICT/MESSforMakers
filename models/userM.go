package models

import (
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID    int
	Name  string
	Email string
}

//get one user (need user ID populated)
func (u *User) getUser(db *sqlx.DB) error {
	return nil
}

//get "count" many users, starting "offset" users from the beginning
func getAllUsers(db *sqlx.DB, count, offset int) ([]User, error) {
	return nil, nil
}

//create user (need user details populated)
func (u *User) createUser(db *sqlx.DB) error {
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

//define Valuers and Scanners for any user related custom types here.
