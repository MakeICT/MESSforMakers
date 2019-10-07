package main

import (
	"fmt"
	"net/http"
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
	config, err := util.InitConfig("config.json")
	if err != nil {
		//No log files have been set up yet, so just dump the error to stdout
		fmt.Print("Cannot parse the configuration file")
		panic(1)
	}

	//TODO newApplication should return an error if the configuration is invalid
	app, err := newApplication(config)
	if err != nil {
		fmt.Printf("Could not start the application: %v", err)
	}
	defer app.logger.Close()

	addr := ":" + strconv.Itoa(app.port)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Addr:         addr,
		ErrorLog:     app.logger.Logger,
		Handler:      app.Router,
	}

	app.logger.Println("Starting Application on :" + strconv.Itoa(app.port))
	app.logger.Fatal(srv.ListenAndServe)

}

// RootHandler is obsolete, will all be handled in controllers.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "you got the root handler")
}
