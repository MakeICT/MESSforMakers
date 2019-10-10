package models

import (
	"github.com/jmoiron/sqlx"
)

//InitDB connects to the database and checks the connection
func InitDB(dataSourceName string) (*sqlx.DB, error) {

	//open returns no error as no connection is made until needed
	db, _ := sqlx.Open("postgres", dataSourceName)
	//force a connection and check that it works.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
