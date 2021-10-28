package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	r := httprouter.New()

	r.NotFound = http.HandlerFunc(app.notFoundErrRespone)
	r.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedErrResponse)

	r.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheck)
	r.HandlerFunc(http.MethodPost, "/v1/users", app.registerUser)
	r.HandlerFunc(http.MethodPost, "/v1/tokens/auth", app.createAuthToken)

	return app.recoverFromPanic(r)
}
