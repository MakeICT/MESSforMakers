package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (a *application) appRouter() {

	//authMiddleware := authenticationMiddleware{a.Session, a.DB}

	//middleware that should be called on every request get added to the chain here
	standardChain := alice.New(a.recoverPanic, a.securityHeaders, a.loggingHandler)
	authenticateChain := alice.New(a.Session.Enable, a.authenticationHandler)

	router := mux.NewRouter()

	fs := http.StripPrefix("/assets/", http.FileServer(noListFileSystem{http.Dir("./assets")}))
	// Ending the URL path with no trailing slash looks like it should be a route to a resource, but it's not. Return 404.
	// Anything that does have the trailing slash may be a real request for an asset. Handle with the custom filesystem
	router.PathPrefix("/assets/").Handler(fs)

	router.Handle("/user", authenticateChain.ThenFunc(a.UserC.Show()))
	router.Handle("/user", authenticateChain.ThenFunc(noRoute("none")))

	//set all the routes here. Uses gorilla/mux so routes can use regex,
	//and following with .Methods() allows for limiting them to only specific HTTP methods
	router.HandleFunc("/", a.StaticC.Root())
	router.HandleFunc("/signup", a.UserC.SignupForm()).Methods("GET")
	router.HandleFunc("/signup", a.UserC.Signup()).Methods("POST")
	router.HandleFunc("/login", a.UserC.LoginForm()).Methods("GET")
	router.HandleFunc("/login", a.UserC.Login()).Methods("POST")
	router.HandleFunc("/logout", a.UserC.Logout()).Methods("POST")
	router.HandleFunc("/user", a.UserC.ShowCurrent()).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}", a.UserC.Show()).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/edit", a.UserC.EditForm()).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}", a.UserC.Edit()).Methods("POST").MatcherFunc(makeMatcher("patch"))
	router.HandleFunc("/user/{id:[0-9]+}", a.UserC.Delete()).Methods("POST").MatcherFunc(makeMatcher("delete"))
	router.HandleFunc("/users", a.UserC.List()).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/ice", noRoute("update ice")).Methods("POST").MatcherFunc(makeMatcher("patch"))
	router.HandleFunc("/user/{id:[0-9]+}/ice", noRoute("delete ice")).Methods("POST").MatcherFunc(makeMatcher("delete"))
	router.HandleFunc("/user/{id:[0-9]+}/uploadwaiver", noRoute("uploadwaiver")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/uploadwaiver", noRoute("save waiver")).Methods("POST")
	router.HandleFunc("/user/{id:[0-9]+}/waiver", noRoute("show waiver")).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/waiver", noRoute("delete waiver")).Methods("POST").MatcherFunc(makeMatcher("delete"))
	router.HandleFunc("/admin", noRoute("admin dashboard")).Methods("GET")

	//TODO: need to implement handlers for 404 and 405, then implement router.NotFoundHandler and router.MethodNotAllowedHandler

	//set the app router. Alice will pass all the requests through the middleware chain first,
	//then to the functions defined above
	a.Router = standardChain.Then(router)
}

func makeMatcher(m string) func(*http.Request, *mux.RouteMatch) bool {
	return func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.FormValue("_method") == m
	}
}

//return a generic error for routes that are in the table but controllers are not ready
func noRoute(msg string) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf("no handler for %s: ", msg))
	})
}

//Prevent directory listings:
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
