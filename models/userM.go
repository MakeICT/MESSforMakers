package models

type UserModel struct {
	Model
}

type User struct {
	ID    int
	Name  string
	Email string
}

//get one user (need user ID populated)
func (um *UserModel) GetUser(int) (*User, error) {
	return nil, nil
}

//get "count" many users, starting "offset" users from the beginning
func (um *UserModel) GetAllUsers(count, offset int) ([]User, error) {
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
func (um *UserModel) CreateUser(*User) error {
	return nil
}

//update user (need user details populated)
func (um *UserModel) UpdateUser(*User) error {
	return nil
}

//delete user (need user ID populated)
func (um *UserModel) DeleteUser(*User) error {
	return nil
}

//define database helper functions here

//define Valuers and Scanners for any user related custom types here.
