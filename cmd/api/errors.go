package main

import (
	"fmt"
	"net/http"
)

func (app *application) sendErrorResponse(w http.ResponseWriter, r *http.Request, status int, msg any) {
	err := app.writeJSON(w, status, payload{"error": msg})
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "url", r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error(), "method", r.Method, "url", r.URL.String())
	msg := "the server encountered a problem and could not process your request"
	app.sendErrorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	msg := "the requested resource could not be found"
	app.sendErrorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.sendErrorResponse(w, r, http.StatusMethodNotAllowed, msg)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.sendErrorResponse(w, r, http.StatusBadRequest, err.Error())
}
