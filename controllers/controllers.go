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
package controllers

import (
	"net/http"

	"github.com/makeict/MESSforMakers/views"
)

type Controller struct{}

//this file for defining methods common to all controllers

func NotImplementedController() Controller {
	return Controller{}
}

func (c *Controller) None(route string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body := "This route has not been implemented yet: " + route

		if err := views.ErrorPage.Index.Render(w, body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}
