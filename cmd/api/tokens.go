package main

import (
	"net/http"
	"time"

	"github.com/mycok/shopit/internal/data"
	"github.com/mycok/shopit/internal/validator"
)

func (app *application) createAuthTokenHandler(rw http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(rw, r, &input)
	if err != nil {
		app.serverErrResponse(rw, r, err)

		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)
	if !v.IsValid() {
		app.failedValidationResponse(rw, r, v.Errors)

		return
	}

	var user data.User

	err = app.repositories.Users.GetByEmail(input.Email, &user)
	if err != nil {
		app.invalidCredentialsResponse(rw, r)

		return
	}

	mathes, err := user.Password.DoesMatch(input.Password)
	if err != nil {
		app.serverErrResponse(rw, r, err)

		return
	}

	if !mathes {
		app.invalidCredentialsResponse(rw, r)

		return
	}

	token, err := app.repositories.Tokens.New(24*time.Hour, user.ID, data.ScopeAuthentication)
	if err != nil {
		app.serverErrResponse(rw, r, err)

		return
	}

	err = app.writeJSON(rw, http.StatusCreated, envelope{
		"auth_token": token,
	}, nil)
	if err != nil {
		app.serverErrResponse(rw, r, err)
	}
}