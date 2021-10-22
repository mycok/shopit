package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	r := httprouter.New()

	r.NotFound = http.HandlerFunc(app.notFoundErrRespone)
	r.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedErrResponse)

	r.HandlerFunc(http.MethodGet, "/healthcheck", app.healthCheckHandler)

	return app.recoverFromPanic(r)
}
