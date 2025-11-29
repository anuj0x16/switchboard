package main

import (
	"fmt"
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, msg any) {
	err := app.writeJSON(w, status, payload{"error": msg})
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "url", r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "url", r.URL.String())
	msg := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	msg := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidation(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
