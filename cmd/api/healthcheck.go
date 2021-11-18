package main

import (
	"net/http"
)

func (app *application) healthCheck(rw http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]interface{}{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(rw, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrResponse(rw, r, err)
	}
}
