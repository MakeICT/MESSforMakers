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
	ID               int       `db:"id"`
	Name             string    `db:"name"`
	Email            string    `db:"username"`
	Password         string    `db:"password"`
	DOB              time.Time `db:"dob"`
	Phone            string    `db:"phone"`
	TextOK           bool      `db:"text_ok"`
	MembershipStatus int       `db:"membership_status_id"`
	MembershipOption int       `db:"membership_option"`
	RBACRole         int       `db:"rbac_role_id"`
}

//Get one user (need user ID populated)
func (um *UserModel) Get(int) (*User, error) {
	return nil, nil
}

//GetAll returns "count" many users, starting "offset" users from the beginning
func (um *UserModel) GetAll(count, page int, sortBy, direction string) ([]User, error) {
	offset := (page - 1) * count //for page 1 the offset should be 0, etc.
	q := um.DB.Rebind("SELECT id, name, username, dob, phone FROM member ORDER BY name LIMIT ? OFFSET ?")
	rows, err := um.DB.Queryx(q, count, offset)
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
func (um *UserModel) Create(u *User) error {
	//TODO calculate membership_expires
	q := um.DB.Rebind(`
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
		u.Name,
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

//define database helper functions here

//define Valuers and Scanners for any user related custom types here.
