/*
 Domesticated Apricot - An open source member and event management platform
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
	"net/http"
	"path/filepath"
)

//generic defaults pages for each controller
type View struct {
	Index Page
	Show  Page
	New   Page
	Edit  Page
}

type Page struct {
	Template *template.Template
	Layout   string
}

func (self *Page) Render(w http.ResponseWriter, data interface{}) error {
	return self.Template.ExecuteTemplate(w, self.Layout, data)
}

func LayoutFiles() []string {
	files, err := filepath.Glob("templates/layouts/*.gohtml")
	if err != nil {
		log.Panic(err)
	}
	return files
}
