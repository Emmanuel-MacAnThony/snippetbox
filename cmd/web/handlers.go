package main

import (
	"errors"
	"fmt"
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
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(response)
		} else {
			app.serverError(response, err)
		}
		return
	}

	data := app.newTemplateData(*request)
	data.Snippet = snippet

	// Use the new render helper.
	app.render(response, http.StatusOK, "view.tmpl.html", data)

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
	if err != nil {
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

	data := app.newTemplateData(*request)
	data.Snippets = snippets

	
	//Use the new render helper.
	app.render(response, http.StatusOK, "home.tmpl.html", data)
	

}
