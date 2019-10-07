package models

import (
	"github.com/jmoiron/sqlx"
)

type UserModel struct {
	DB *sqlx.DB
}

type User struct {
	ID    int
	Name  string
	Email string
}

//get one user (need user ID populated)
func (um *UserModel) Get(int) (*User, error) {
	return nil, nil
}

//get "count" many users, starting "offset" users from the beginning
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

//create user (need user details populated)
func (um *UserModel) Create(*User) error {
	return nil
}

//update user (need user details populated)
func (um *UserModel) Update(*User) error {
	return nil
}

//delete user (need user ID populated)
func (um *UserModel) Delete(*User) error {
	return nil
}

//define database helper functions here

//define Valuers and Scanners for any user related custom types here.
