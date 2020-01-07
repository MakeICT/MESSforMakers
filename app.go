package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/makeict/MESSforMakers/controllers"
	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/util"
)

// database connection, cookie store, etc..
type application struct {
	Logger  *util.Logger
	DB      *sqlx.DB
	Router  http.Handler
	Config  *util.Config
	UserC   controllers.UserController
	StaticC controllers.StaticController
	ErrorC  controllers.ErrorController
	Session *sessions.Session
	port    int
}

func newApplication(config *util.Config) (*application, error) {

	//Set up a logger middleware
	logger, err := util.NewLogger(config.Logger.LogFile, config.Logger.DumpRequest, util.DEBUG)
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %v\n", err)
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
		return nil, fmt.Errorf("error initializing database: %v\n", err)
	}
	session := sessions.New([]byte(config.App.SessionKey))
	session.Lifetime = 12 * time.Hour
	app.Session = session
	app.DB = db
	app.port = config.App.Port

	um := &models.UserModel{DB: app.DB, HashCost: 12}

	if err := app.UserC.Initialize(app.Config, um, app.Logger, app.Session); err != nil {
		app.Logger.Fatalf("failed to initialize user controller: %v\n", err)
	}

	if err := app.StaticC.Initialize(app.Config, um, app.Logger, app.Session); err != nil {
		app.Logger.Fatalf("failed to initialize controller for static routes: %v\n", err)
	}

	if err := app.ErrorC.Initialize(app.Config, um, app.Logger, app.Session); err != nil {
		app.Logger.Fatalf("failed to initialize controller for error routes: %v\n", err)
	}

	//initialize all the routes
	app.appRouter()

	return &app, nil
}
