package main

import "net/http"

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	data := payload{"status": "available"}

	err := app.writeJSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
