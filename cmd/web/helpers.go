package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) newTemplateData(request *http.Request) *templateData{
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) render(response http.ResponseWriter, status int, page string, data *templateData) {

	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(response, err)
		return
	}

	// initialize new buffer
	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(response, err)
		return
	}

	response.WriteHeader(status)

	buf.WriteTo(response)
}

func (app *application) serverError(response http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(response http.ResponseWriter, status int) {
	http.Error(response, http.StatusText(status), status)
}

func (app *application) notFound(response http.ResponseWriter) {
	app.clientError(response, http.StatusNotFound)
}
