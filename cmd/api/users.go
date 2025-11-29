package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/anuj0x16/switchboard/internal/data"
	"github.com/anuj0x16/switchboard/internal/validator"
	"github.com/golang-jwt/jwt/v5"
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

func (app *application) authLoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentials(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	match, err := user.Password.Match(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !match {
		app.invalidCredentials(w, r)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Subject:   user.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString([]byte(app.config.jwtSecret))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, payload{"authentication_token": tokenString})
	if err != nil {
		app.serverError(w, r, err)
	}
}
