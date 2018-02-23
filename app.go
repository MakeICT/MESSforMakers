package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

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
	port        int
}

func newApplication(config *Config) *application {

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

	//Initialize the cookie store
	cs := session.NewCookieStore("mess-data")

	//set up authentication
	authMiddleware := authenticationMiddleware{cs, db}

	//Set up a logger middleware
	logger, err := util.NewLogger()
	if err != nil {
		fmt.Printf("Error creating logger :: %v", err)
		panic(1)
	}
	loggingMiddleware := loggingMiddleware{config.Logger.DumpRequest, logger}

	//middleware that should be called on every request get added to the chain here
	commonHandlers := alice.New(loggingMiddleware.loggingHandler, authMiddleware.authenticationHandler)

	// app needs to be created and have the DB initialized so that appRouter can pass the connection pool to the controllers
	app := application{
		cookieStore: cs,
		logger:      logger,
		DB:          db,
		port:        config.App.Port,
	}

	//initialize all the routes
	app.appRouter(commonHandlers)

	return &app
}

func (a *application) appRouter(c alice.Chain) {

	router := mux.NewRouter()

	//create the controllers so they are more readable in the routes table
	userC := controllers.User(a.DB, a.cookieStore)
	NIC := controllers.NotImplementedController()
	staticC := controllers.StaticController()

	//set all the routes here. Uses gorilla/mux so routes can use regex,
	//and following with .Methods() allows for limiting them to only specific HTTP methods
	// Routes that utilize the same controller and switch inside the controller on method type can use the same table entry.
	router.HandleFunc("/", staticC.Root()).Methods("GET")
	router.HandleFunc("/join", staticC.Join()).Methods("GET")
	router.HandleFunc("/reserve", staticC.Reservations()).Methods("GET")
	router.HandleFunc("/signup", userC.Create()).Methods("GET", "POST")
	router.HandleFunc("/login", userC.Login()).Methods("GET", "POST")
	router.HandleFunc("/user/{id:[0-9]+}", userC.Show()).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/edit", userC.Edit()).Methods("GET", "POST")
	router.HandleFunc("/user/{id:[0-9]+}", NIC.None("delete user")).Methods("DELETE")
	router.HandleFunc("/users", userC.Index()).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/ice", NIC.None("update ice")).Methods("PATCH")
	router.HandleFunc("/user/{id:[0-9]+}/ice", NIC.None("delete ice")).Methods("DELETE")
	router.HandleFunc("/user/{id:[0-9]+}/uploadwaiver", NIC.None("uploadwaiver")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/uploadwaiver", NIC.None("save waiver")).Methods("POST")
	router.HandleFunc("/user/{id:[0-9]+}/waiver", NIC.None("show waiver")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/waiver", NIC.None("delete waiver")).Methods("DELETE")
	router.HandleFunc("/admin", NIC.None("admin dashboard")).Methods("GET")

	//TODO: need better static file serving to prevent directory browsing
	// see https://groups.google.com/forum/#!topic/golang-nuts/bStLPdIVM6w
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//TODO: need to implement handlers for 404 and 405, the implement router.NotFoundHandler and router.MethodNotAllowedHandler

	//set the app router. Alice will pass all the requests through the middleware chain first,
	//then to the functions defined above
	a.Router = c.Then(router)
}

//Middleware
type loggingMiddleware struct {
	dumpRequest bool
	logger      *util.Logger
}

func (l *loggingMiddleware) loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		if l.dumpRequest {
			if reqDump, err := httputil.DumpRequest(r, true); err == nil {
				l.logger.Printf("Recieved Request:\n%s\n", reqDump)
			}
		}

		h.ServeHTTP(w, r)

		t2 := time.Now()

		l.logger.Printf("[%s] %s %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	})
}

type authenticationMiddleware struct {
	CookieStore *session.CookieStore
	DB          *sqlx.DB
}

type contextKey int

const contextkeyUser contextKey = 1

func (a *authenticationMiddleware) authenticationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//check if cookie exists
		s, err := a.CookieStore.Store.Get(r, "mess-data")

		if err == nil && !s.IsNew { //if the session is not new, then it can look for the user ID and auth key

			id := s.Values["userID"].(string)
			auth := s.Values["authToken"].(string)

			//if it exists, check user ID and token for validity.
			if id != "" && auth != "" {

				user, err := models.SessionLookup(a.DB, id, auth)

				//if username and token are valid, store the username, id, and (later) rbac role in the context for use by later handlers
				if err == nil {
					user.Authorized = true
					// context key should be a custom type and a const NOT a string
					ctx := context.WithValue(r.Context(), contextkeyUser, user)

					// TODO if the user is authorized, the LastSeenTime of the session should be updated

					h.ServeHTTP(w, r.WithContext(ctx))
					return
				} else if err == models.ErrNotAuthorized {
					// user is either not allowed to log in, or the access token is expired.
					// is there a distinction? Always redirect to login form after setting a flash?
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		h.ServeHTTP(w, r)

	})
}
