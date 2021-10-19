package main

import (
	"net/http"
	"fmt"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: http.DefaultServeMux,
	}

	app.logger.LogInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env": app.config.env,
	})

	return srv.ListenAndServe()
}