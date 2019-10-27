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
	Session *sessions.Session
	port    int
}

func newApplication(config *util.Config) (*application, error) {

	//Set up a logger middleware
	logger, err := util.NewLogger("makeict.log", config.Logger.DumpRequest, util.DEBUG) //KNOWN BUG, first arg is ignored
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
	session := sessions.New([]byte("H97LwY3g5X5W0AJjdw4yEIZCIasiU2FRg"))
	session.Lifetime = 12 * time.Hour
	app.Session = session
	app.DB = db
	app.port = config.App.Port

	if err := app.UserC.Initialize(app.Config, &models.UserModel{DB: app.DB}, app.Logger, app.Session); err != nil {
		app.Logger.Fatalf("Failed to initialize user controller: %v", err)
	}

	if err := app.StaticC.Initialize(app.Config, &models.UserModel{DB: app.DB}, app.Logger, app.Session); err != nil {
		app.Logger.Fatalf("Failed to initialize controller for static routes: %v", err)
	}

	//initialize all the routes
	app.appRouter()

	return &app, nil
}
