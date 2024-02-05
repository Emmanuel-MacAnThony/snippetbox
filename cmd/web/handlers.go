package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	
	"github.com/Emmanuel-MacAnThony/snippetbox/internal/models"
	"github.com/Emmanuel-MacAnThony/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type SnippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	validator.Validator
}

// snippet view handler function
func (app *application) snippetView(response http.ResponseWriter, request *http.Request) {

	params := httprouter.ParamsFromContext(request.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

	data := app.newTemplateData(request)
	data.Snippet = snippet

	// Use the new render helper.
	app.render(response, http.StatusOK, "view.tmpl.html", data)

}

// snippet create handler function
func (app *application) snippetCreatePost(response http.ResponseWriter, request *http.Request) {
	// parse form to extract data in request body to postform map
	err := request.ParseForm()
	if err != nil {
		app.clientError(response, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(request.PostForm.Get("expires"))
	if err != nil {
		app.clientError(response, http.StatusBadRequest)
		return
	}

	form := SnippetCreateForm{
		Title:       request.PostForm.Get("title"),
		Content:     request.PostForm.Get("content"),
		Expires:     expires,
		
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid(){
		data := app.newTemplateData(request)
		data.Form = form
		app.render(response, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(response, err)
		return
	}

	http.Redirect(response, request, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreate(response http.ResponseWriter, request *http.Request) {
	data := app.newTemplateData(request)
	data.Form = SnippetCreateForm{
		Expires: 365,
	}
	app.render(response, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) home(response http.ResponseWriter, request *http.Request) {
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(response, err)
		return
	}

	data := app.newTemplateData(request)
	data.Snippets = snippets

	//Use the new render helper.
	app.render(response, http.StatusOK, "home.tmpl.html", data)
}
