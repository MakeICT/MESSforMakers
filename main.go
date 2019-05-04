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
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/makeict/MESSforMakers/util"

	"html/template"
)

const appPort = "3000"

func main() {

	// Initialize the configuration struct from the configuration file
	config, err := InitConfig("config.json")
	if err != nil {
		fmt.Print("Cannot parse the configuration file")
		panic(1)
	}

	// create the app with user-defined settings
	app := newApplication(config)

	// make sure the logger releases its resources if the server shuts down.
	defer app.logger.Close()
	app.logger.Warn("this is a test warning")
	app.logger.Println("Starting Application")
	app.logger.Fatal(http.ListenAndServe("localhost:"+appPort, app.Router))

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

// RootHandler is obsolete, will all be handled in controllers.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	var name string

	nameCookie, err := r.Cookie("name")
	if err == nil {
		name = nameCookie.Value
	} else if err != http.ErrNoCookie {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := template.Must(template.New("root").Parse(`
	<!doctype html>
	<html>
	<head><title>greeter</title></head>
	<body>
	{{if .}}
	  <h1>Hi {{.}}!</h1>
	{{else}}
	  <h1>You got the root handler !</h1>
	  <h2>Who are you?</h2>
	  <form method="POST" action="/">
	    <input type="text" placeholder="Your name" name="name">
	    <input type="submit" value="Set Name" id="submit-button">
	  </form>
	{{end}}
	</body>
	</html>
	`))
	if err := t.Execute(w, name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "name",
		Value:    r.FormValue("name"),
		HttpOnly: true,
		Expires:  time.Now().Add(24 * 14 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}
