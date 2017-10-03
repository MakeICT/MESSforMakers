/*
 MESS for Makers - An open source member and event management platform
    Copyright (C) 2017  Sam Schurter

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/makeict/MESSforMakers/controllers"
	"log"
	"net/http"
	"time"
)

type App struct {
	*mux.Router
}

func main() {
	app := &App{mux.NewRouter()}
	commonHandlers := alice.New(loggingHandler)
	userC := controllers.User
	app.Handle("/", commonHandlers.ThenFunc(RootHandler))
	app.Handle("/user", commonHandlers.ThenFunc(userC.Index))
	log.Fatal(http.ListenAndServe(":8080", app))
}

//middleware
func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		h.ServeHTTP(w, r)
		t2 := time.Now()
		fmt.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	})
}

//obsolete, will all be handled in controllers
func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "you got the root handler")
}
