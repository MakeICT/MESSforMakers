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

const appPort = "8080"

func main() {
	config, err := InitConfig("config.json")
	if err != nil {
		fmt.Print("Cannot parse the configuration file")
		panic(1)
	}
	app := newApplication(config)
	defer app.logger.Close()

	app.logger.Println("Starting Application")
	app.logger.Fatal(http.ListenAndServe(":"+strconv.Itoa(app.port), app.Router))

}
