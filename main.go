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
	"strconv"
	"time"

	"github.com/makeict/MESSforMakers/util"
)

//Main reads the configuration immediately and dies if it can't be read.
//There is no default configuration and no command line flags.
//All options are contained in the config.json file
func main() {
	//Read the configuration and die if it can't be read.
	//This does NOT guarantee that sensible options have been set, only that the file can be read.
	config, err := InitConfig("config.json")
	if err != nil {
		//No log files have been set up yet, so just dump the error to stdout
		fmt.Print("Cannot parse the configuration file")
		panic(1)
	}

	//TODO newApplication should return an error if the configuration is invalid
	app := newApplication(config)
	defer app.logger.Close()

	app.logger.Println("Starting Application on :" + strconv.Itoa(app.port))
	app.logger.Fatal(http.ListenAndServe(":"+strconv.Itoa(app.port), app.Router))

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
	fmt.Fprintln(w, "you got the root handler")
}
