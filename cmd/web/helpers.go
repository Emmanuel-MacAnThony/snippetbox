package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *application) newTemplateData(request *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// Add the flash message to the template data, if one exists.
		Flash: app.sessionManager.PopString(request.Context(), "flash"),
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

func (app *application) decodePostForm(request *http.Request, dst any) error {

	err := request.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, request.PostForm)
	if err != nil {

		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}
	return nil

}
