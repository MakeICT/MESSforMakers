package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/makeict/MESSforMakers/controllers"
	"github.com/makeict/MESSforMakers/session"
	"github.com/makeict/MESSforMakers/util"
)

// database connection, cookie store, etc..
type application struct {
	cookieStore    *session.CookieStore
	commonHandlers alice.Chain
	logger         *util.Logger
}

func newApplication() *application {
	logger, err := util.NewLogger()
	if err != nil {
		fmt.Printf("Error creating logger :: %v", err)
		panic(1)
	}
	loggingMiddleware := loggingMiddleware{true, logger}

	return &application{
		cookieStore:    session.NewCookieStore("mess-data"),
		commonHandlers: alice.New(loggingMiddleware.loggingHandler),
		logger:         logger,
	}
}

func (a *application) appRouter() http.Handler {
	router := mux.NewRouter()
	userC := controllers.User

	router.HandleFunc("/", RootHandler)
	router.HandleFunc("/user", userC.Index)

	return a.commonHandlers.Then(router)
}
