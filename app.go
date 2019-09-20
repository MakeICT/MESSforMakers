package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/makeict/MESSforMakers/controllers"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
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

	//initialize all the routes
	app.appRouter()

	return &app, nil
}

func noRoute(msg string) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf("no route for %s: ", msg))
	})
}

// A controller is a struct that stores pointers to all the necessary things
// from the application level. Handlers are defined on that struct.

func (a *application) appRouter() {

	//middleware that should be called on every request get added to the chain here
	c := alice.New(a.recoverPanic, a.loggingHandler)

	router := mux.NewRouter()

	fs := http.StripPrefix("/assets/", http.FileServer(noListFileSystem{http.Dir("./assets")}))
	// Ending the URL path with no trailing slash looks like it should be a route to a resource, but it's not. Return 404.
	// Anything that does have the trailing slash may be a real request for an asset. Handle with the custom filesystem
	router.PathPrefix("/assets/").Handler(fs)

	userC := controllers.UserController{}
	if err := userC.Initialize(a.Config, a.CookieStore, a.DB, a.Logger); err != nil {
		a.Logger.Fatalf("Failed to initialize user controller: %v", err)
	}

	staticC := controllers.StaticController{}
	if err := staticC.Initialize(a.Config, a.CookieStore, a.DB, a.Logger); err != nil {
		a.Logger.Fatalf("Failed to initialize controller for static routes: %v", err)
	}

	//set all the routes here. Uses gorilla/mux so routes can use regex,
	//and following with .Methods() allows for limiting them to only specific HTTP methods
	router.HandleFunc("/", staticC.Root())
	router.HandleFunc("/signup", noRoute("signup")).Methods("GET")
	router.HandleFunc("/login", noRoute("login")).Methods("GET")
	router.HandleFunc("/login", noRoute("login processor")).Methods("POST")
	router.HandleFunc("/user", noRoute("save new user to db")).Methods("POST")
	router.HandleFunc("/user/{id:[0-9]+}", noRoute("show specific user")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/edit", noRoute("form to edit user")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}", noRoute("save user update to db")).Methods("PATCH")
	router.HandleFunc("/user/{id:[0-9]+}", noRoute("delete user")).Methods("DELETE")
	router.HandleFunc("/users", noRoute("users")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/ice", noRoute("update ice")).Methods("PATCH")
	router.HandleFunc("/user/{id:[0-9]+}/ice", noRoute("delete ice")).Methods("DELETE")
	router.HandleFunc("/user/{id:[0-9]+}/uploadwaiver", noRoute("uploadwaiver")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/uploadwaiver", noRoute("save waiver")).Methods("POST")
	router.HandleFunc("/user/{id:[0-9]+}/waiver", noRoute("show waiver")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/waiver", noRoute("delete waiver")).Methods("DELETE")

	//TODO: need to implement handlers for 404 and 405, then implement router.NotFoundHandler and router.MethodNotAllowedHandler

	//set the app router. Alice will pass all the requests through the middleware chain first,
	//then to the functions defined above
	a.Router = c.Then(router)
}

type noListFileSystem struct {
	fs http.FileSystem
}

func (nfs noListFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		if _, err := nfs.fs.Open(strings.TrimSuffix(path, "/") + "/index.html"); err != nil {
			return nil, err
		}
	}

	return f, nil
}
