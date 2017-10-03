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
package views

import (
	"html/template"
	"log"
	"path/filepath"
)

type UserView struct {
	View
	//add custom pages for a controller here
}

var User UserView

func UserFiles() []string {
	files, err := filepath.Glob("templates/user/includes/*.gohtml")
	if err != nil {
		log.Panic(err)
	}
	files = append(files, LayoutFiles()...)
	return files
}

func init() {
	indexFiles := append(UserFiles(), "templates/user/index.gohtml")
	User.Index = Page{
		Template: template.Must(template.New("index").ParseFiles(indexFiles...)),
		Layout:   "index",
	}

	//parse other needed templates here, for each page

}
