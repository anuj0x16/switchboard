package main

import (
	"errors"
	"net/http"

	"github.com/anuj0x16/switchboard/internal/data"
	"github.com/anuj0x16/switchboard/internal/validator"
)

func (app *application) authRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &data.User{
		Email: input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "already in use")
			app.failedValidation(w, r, v.Errors)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, payload{"user": user})
	if err != nil {
		app.serverError(w, r, err)
	}
}
