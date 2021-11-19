package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", app.config.port),
		Handler:     app.routes(),
		ErrorLog:    log.New(app.logger, "", 0),
		IdleTimeout: time.Minute,
	}

	app.server = srv

	// Handle graceful shutdown of the server.
	shutdownErrChan := make(chan error)
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)

		s := <-quitChan

		app.logger.LogInfo("server shutting down", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call Shutdown() on our server, passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because of a problem closing the listeners, or
		// because the shutdown didn't complete before the 5-second context deadline is hit).
		// We relay this return value to the shutdownError channel.
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownErrChan <- err
		}

		app.logger.LogInfo("completing all background tasks", map[string]string{
			"addr": srv.Addr,
		})

		// Wait for all running goroutines to return / exit.
		app.wg.Wait()
		close(shutdownErrChan)
	}()

	app.logger.LogInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrChan
	if err != nil {
		return err
	}

	app.logger.LogInfo("server stopped", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
