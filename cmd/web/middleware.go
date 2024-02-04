package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		response.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		response.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		response.Header().Set("X-Content-Type-Options", "nosniff")
		response.Header().Set("X-Frame-Options", "deny")
		response.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(response, request)

	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", request.RemoteAddr, request.Proto, request.Method, request.URL.RequestURI())
		next.ServeHTTP(response, request)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler{

	return http.HandlerFunc(func (response http.ResponseWriter, request *http.Request) {
			defer func() {
				if err := recover() ; err != nil{
					response.Header().Set("Connection", "close")
					app.serverError(response, fmt.Errorf("%s", err))
				}
			}()
			next.ServeHTTP(response, request)
		})
}

