package models

import (
	"github.com/jmoiron/sqlx"
)

//InitDB connects to the database and checks the connection
func InitDB(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
