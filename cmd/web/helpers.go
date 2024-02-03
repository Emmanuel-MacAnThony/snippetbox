package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(response http.ResponseWriter, err error){
	trace := fmt.Sprintf("%s\n%s",err.Error(),debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(response http.ResponseWriter, status int){
	http.Error(response, http.StatusText(status), status)
}

func (app *application) notFound(response http.ResponseWriter){
	app.clientError(response, http.StatusNotFound)
}