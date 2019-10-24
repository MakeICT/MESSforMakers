package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
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
	CookieStore *sessions.CookieStore
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
