package main

import (
	"errors"
	"fmt"
	//"html/template"
	"net/http"
	"strconv"

	"github.com/Emmanuel-MacAnThony/snippetbox/internal/models"
)

// snippet view handler function
func (app *application) snippetView(response http.ResponseWriter, request *http.Request) {

	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(response)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil{
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(response)
		}else {
			app.serverError(response, err)
			}
			return
	}
	fmt.Fprintf(response, "%+v", snippet)
}

// snippet create handler function
func (app *application) snippetCreate(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		response.Header().Set("Allow", http.MethodPost)
		app.clientError(response, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil{
		app.serverError(response, err)
	}

	http.Redirect(response, request, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) home(response http.ResponseWriter, request *http.Request) {

	if request.URL.Path != "/" {
		http.NotFound(response, request)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(response, err)
		return
	}

	for _, snippet := range snippets{
		fmt.Fprintf(response, "%+v\n", snippet)
	}

	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
		
	// }

	// // parse template
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.errorLog.Println(err.Error())
	// 	app.serverError(response, err)
	// 	return
	// }
	// err = ts.ExecuteTemplate(response, "base", nil)

	// if err != nil {
	// 	app.errorLog.Println(err.Error())
	// 	app.serverError(response, err)
	// 	return
	// }

	
}
