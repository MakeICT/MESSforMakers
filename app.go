package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"

	"github.com/makeict/MESSforMakers/controllers"
	"github.com/makeict/MESSforMakers/models"
	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/util"
)

// database connection, cookie store, etc..
type application struct {
	cookieStore *session.CookieStore
	logger      *util.Logger
	DB          *sqlx.DB
	Router      http.Handler
}

func newApplication(config *Config) *application {

	//Set up a logger middleware
	logger, err := util.NewLogger()
	if err != nil {
		fmt.Printf("Error creating logger :: %v", err)
		panic(1)
	}
	loggingMiddleware := loggingMiddleware{true, logger} // currently defaults dump_request to true, needs to be a config option

	//middlewware that should be called on every request get added to the chain here
	commonHandlers := alice.New(loggingMiddleware.loggingHandler)

	//set up the database
	db, err := models.InitDB(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database,
	))
	if err != nil {
		fmt.Printf("Error initializing database :: %v", err)
		panic(1)
	}

	// app needs to be created and have the DB initialized so that appRouter can pass the connection pool to the controllers
	app := application{
		cookieStore: session.NewCookieStore("mess-data"),
		logger:      logger,
		DB:          db,
	}

	//initialize all the routes
	app.appRouter(commonHandlers)

	return &app
}

func (a *application) appRouter(c alice.Chain) {

	router := mux.NewRouter()

	//declare all the controllers so they are more readable in the routes table
	userC := controllers.User(a.DB)
	NIC := controllers.NotImplementedController()

	//set all the routes here. Uses gorilla/mux so routes can use regex,
	//and following with .Methods() allows for limiting them to only specific HTTP methods
	router.HandleFunc("/", RootHandler)
	router.HandleFunc("/signup", NIC.None("signup page")).Methods("GET")
	router.HandleFunc("/login", NIC.None("login page")).Methods("GET")
	router.HandleFunc("/login", NIC.None("login processor")).Methods("POST")
	router.HandleFunc("/user", NIC.None("save new user to db")).Methods("POST")
	router.HandleFunc("/user/{id:[0-9]+}", NIC.None("show specific user")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/edit", NIC.None("form to edit user")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}", NIC.None("save user update to db")).Methods("PATCH")
	router.HandleFunc("/user/{id:[0-9]+}", NIC.None("delete user")).Methods("DELETE")
	router.HandleFunc("/users", userC.Index()).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/ice", NIC.None("update ice")).Methods("PATCH")
	router.HandleFunc("/user/{id:[0-9]+}/ice", NIC.None("delete ice")).Methods("DELETE")
	router.HandleFunc("/user/{id:[0-9]+}/uploadwaiver", NIC.None("uploadwaiver")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/uploadwaiver", NIC.None("save waiver")).Methods("POST")
	router.HandleFunc("/user/{id:[0-9]+}/waiver", NIC.None("show waiver")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/waiver", NIC.None("delete waiver")).Methods("DELETE")

	//TODO: need better static file serving to prevent directory browsing
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	//TODO: need to implement handlers for 404 and 405, the implement router.NotFoundHandler and router.MethodNotAllowedHandler

	//set the app router. Alice will pass all the requests through the middleware chain first,
	//then to the functions defined above
	a.Router = c.Then(router)
}
