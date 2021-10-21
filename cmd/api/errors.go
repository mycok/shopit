package main

import (
	"fmt"
	"net/http"
)

func (app *application) serverErrorResponse(rw http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errResponse(rw, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundError(rw http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errResponse(rw, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(rw http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errResponse(rw, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(rw http.ResponseWriter, r *http.Request, err error) {
	app.errResponse(rw, r, http.StatusBadRequest, err.Error())
}

func (app *application) logError(r *http.Request, err error) {
	app.logger.LogError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errResponse(rw http.ResponseWriter, r *http.Request, statusCode int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(rw, statusCode, env, nil)
	if err != nil {
		app.logError(r, err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
