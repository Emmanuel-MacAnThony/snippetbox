package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// snippet view handler function
func (app *application) snippetView(response http.ResponseWriter, request *http.Request) {

	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(response)
		return
	}
	fmt.Fprintf(response, "Display a specific snippet with ID %d...", id)
}

// snippet create handler function
func (app *application) snippetCreate(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.Header().Set("Allow", http.MethodPost)
		app.clientError(response, http.StatusMethodNotAllowed)
		return
	}
	response.Write([]byte("Create a new snippet"))
}

func (app *application) home(response http.ResponseWriter, request *http.Request) {

	if request.URL.Path != "/" {
		http.NotFound(response, request)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		
	}

	// parse template
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(response, err)
		return
	}
	err = ts.ExecuteTemplate(response, "base", nil)

	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(response, err)
		return
	}

}
