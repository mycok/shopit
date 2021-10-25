package main

import (
	"fmt"
	"net/http"
)

func (app *application) serverErrResponse(rw http.ResponseWriter, r *http.Request, err error) {
	app.logErr(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errResponse(rw, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundErrRespone(rw http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errResponse(rw, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedErrResponse(rw http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errResponse(rw, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestErrResponse(rw http.ResponseWriter, r *http.Request, err error) {
	app.errResponse(rw, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(rw http.ResponseWriter, r *http.Request, errs map[string]string) {
	app.errResponse(rw, r, http.StatusUnprocessableEntity, errs)
}

func (app *application) logErr(r *http.Request, err error) {
	app.logger.LogError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errResponse(rw http.ResponseWriter, r *http.Request, statusCode int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(rw, statusCode, env, nil)
	if err != nil {
		app.logErr(r, err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
