package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/mycok/shopit/internal/data"
	"github.com/mycok/shopit/internal/validator"
)

func (app *application) registerUser(rw http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(rw, r, &input)
	if err != nil {
		app.badRequestErrResponse(rw, r, err)

		return
	}

	user := &data.User{
		Username:  input.Username,
		Email:     input.Email,
		Version:   version,
		CreatedAt: time.Now().UTC(),
	}

	v := validator.New()
	data.ValidatePassword(v, input.Password)
	if !v.IsValid() {
		app.failedValidationResponse(rw, r, v.Errors)

		return
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrResponse(rw, r, err)

		return
	}

	if user.Validate(v); !v.IsValid() {
		app.failedValidationResponse(rw, r, v.Errors)

		return
	}

	_id, err := app.repositories.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateKey):
			app.badRequestErrResponse(rw, r, err)
		default:
			app.serverErrResponse(rw, r, err)
		}

		return
	}

	token, err := app.repositories.Tokens.New(3*24*time.Hour, *_id, data.ScopeActivation)
	if err != nil {
		app.serverErrResponse(rw, r, err)
	}

	app.runInBackground(func() {
		data := map[string]interface{}{
			"activationToken": token.PlainText,
			"userID":          token.UserID,
		}

		err := app.mailer.Send(data, "user_welcome.go.tmpl", user.Email)
		if err != nil {
			// If there is an error sending the email then we use the
			// app.logger.PrintError() helper to manage it, instead of the
			// app.serverErrorResponse() helper. This is because the email sending functionality
			// runs in a background goroutine which means the handler may return before the email sending
			// goroutine returns. This makes writing a JSON error response redundant since the request-response
			// cycle would already have completed.
			app.logger.LogError(err, nil)
		}
	})

	// Send the client a 202 Accepted status code to indicate that the request has been
	// accepted for processing but the processing has not yet been completed to cater for cases where the handler
	// returns before the email sending functionality is complete.
	err = app.writeJSON(rw, http.StatusAccepted, envelope{
		"id": *_id,
	}, nil)
	if err != nil {
		app.serverErrResponse(rw, r, err)
	}
}
