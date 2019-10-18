package main

import (
	"crypto/tls"
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
		fmt.Printf("Could not create the application: %v", err)
	}
	defer app.Logger.Close()

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	tlsSrv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Addr:         fmt.Sprintf("%s:%d", config.App.Host, config.App.Port),
		ErrorLog:     app.Logger.Logger,
		Handler:      app.Router,
		TLSConfig:    tlsConfig,
	}

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + req.Host + req.URL.String()
			http.Redirect(w, req, url, http.StatusMovedPermanently)
		}),
	}

	go func() { app.Logger.Fatal(srv.ListenAndServe()) }()

	app.Logger.Println("Starting Application on :" + strconv.Itoa(app.port))
	app.Logger.Fatal(tlsSrv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))

}
