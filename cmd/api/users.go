package main

import (
	"net/http"
	"time"

	"github.com/mycok/shopit/internal/data"
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
		IsActive: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrResponse(rw, r, err)

		return
	}

	_id, err := app.repositories.Users.Insert(user)
	if err != nil {
		app.serverErrResponse(rw, r, err)

		return
	}

	err = app.writeJSON(rw, http.StatusCreated, envelope{
		"id": _id,
	}, nil)
	if err != nil {
		app.serverErrResponse(rw, r, err)
	}
}
