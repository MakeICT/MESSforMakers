package main

import (
	"fmt"
	"net/http"

	"github.com/makeict/MESSforMakers/controllers"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/util"
)

// database connection, cookie store, etc..
type application struct {
	CookieStore *session.CookieStore
	Logger      *util.Logger
	DB          *sqlx.DB
	Router      http.Handler
	port        int
	Config      *util.Config
	UserC       controllers.UserController
	StaticC     controllers.StaticController
}

func newApplication(config *util.Config) (*application, error) {

	//Set up a logger middleware
	logger, err := util.NewLogger(config.Logger.DumpRequest)
	if err != nil {
		return nil, fmt.Errorf("Error creating logger :: %v", err)
	}

	app := application{Logger: logger, Config: config}

	//set up the database
	db, err := models.InitDB(fmt.Sprintf(
		"sslmode=%s user=%s password=%s host=%s port=%d dbname=%s",
		config.Database.SSL,
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database,
	))
	if err != nil {
		return nil, fmt.Errorf("Error initializing database :: %v", err)
	}

	app.CookieStore = session.NewCookieStore("mess-data")
	app.DB = db
	app.port = config.App.Port

	if err := app.UserC.Initialize(app.Config, app.CookieStore, &models.UserModel{DB: app.DB}, app.Logger); err != nil {
		app.Logger.Fatalf("Failed to initialize user controller: %v", err)
	}

	if err := app.StaticC.Initialize(app.Config, app.CookieStore, &models.UserModel{DB: app.DB}, app.Logger); err != nil {
		app.Logger.Fatalf("Failed to initialize controller for static routes: %v", err)
	}

	//initialize all the routes
	app.appRouter()

	return &app, nil
}
