package main

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/mycok/shopit/internal/data"
	"github.com/mycok/shopit/internal/validator"
)

func (app *application) registerUserHandler(rw http.ResponseWriter, r *http.Request) {
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
		ID:        uuid.New().String(),
		Username:  input.Username,
		Email:     input.Email,
		IsActive:  false,
		IsSeller:  false,
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

		return
	}

	app.runInBackground(func() {
		var serverAddr string

		if app.server.Addr == ":4000" {
			serverAddr = "http://localhost" + app.server.Addr + "/v1/users/activate/" + token.PlainText
		} else {
			serverAddr = app.server.Addr + "/v1/users/activate/" + token.PlainText
		}

		url, err := url.Parse(serverAddr)
		if err != nil {
			app.logger.LogError(err, nil)
		}

		data := map[string]interface{}{
			"activationToken": token.PlainText,
			"activationLink":  url.String(),
		}

		err = app.mailer.Send(data, "user_welcome.go.tmpl", user.Email)
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
		"user_id": _id,
	}, nil)
	if err != nil {
		app.serverErrResponse(rw, r, err)
	}
}

func (app *application) activateUserHandler(rw http.ResponseWriter, r *http.Request) {
	plainTextToken := app.readParam(r, "token")

	v := validator.New()
	if data.ValidatePlainTextToken(v, plainTextToken); !v.IsValid() {
		app.failedValidationResponse(rw, r, v.Errors)

		return
	}

	var token data.Token

	err := app.repositories.Tokens.Get(plainTextToken, data.ScopeActivation, &token)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			// TODO: handle expired token scenario.
			app.badRequestErrResponse(rw, r, data.ErrInvalidOrExpiredToken)
		default:
			app.serverErrResponse(rw, r, err)
		}

		return
	}

	var user data.User

	// TODO: change the GetByID call into a middleware for a group of user related handlers.
	// This change is more suitable for use with JWT tokens.
	err = app.repositories.Users.GetByID(token.UserID, &user)
	if err != nil {
		app.serverErrResponse(rw, r, err)

		return
	}

	updateData := struct {
		IsActive bool
		Version  int64
	}{
		IsActive: true,
		Version:  user.Version + version,
	}

	updateRes, err := app.repositories.Users.Update(user.ID, updateData)
	if err != nil {
		app.serverErrResponse(rw, r, err)

		return
	}

	// Delete all activation tokens that belong to the current user.This mitigates the user
	// from trying to activate more than once.
	_, err = app.repositories.Tokens.DeleteAllForUser(user.ID, data.ScopeActivation)
	if err != nil {
		app.serverErrResponse(rw, r, err)

		return
	}

	err = app.writeJSON(rw, http.StatusOK, envelope{
		"status": updateRes.ModifiedCount,
	}, nil)
	if err != nil {
		app.serverErrResponse(rw, r, err)
	}
}
