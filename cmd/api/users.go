package main

import (
	"net/http"
	"errors"
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
		CreatedAt: time.Now(),
		IsActive:  false,
		IsSeller:  false,
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
		switch  {
		case errors.Is(err, data.DuplicateKeyErr):
			app.badRequestErrResponse(rw, r, err)
		default:
			app.serverErrResponse(rw, r, err)
		}

		return
	}

	err = app.writeJSON(rw, http.StatusCreated, envelope{
		"id": _id,
	}, nil)
	if err != nil {
		app.serverErrResponse(rw, r, err)
	}
}
