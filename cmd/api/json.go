package main

import (
	"encoding/json"
	"net/http"
)

type payload map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data payload) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(append(b, '\n'))

	return nil
}
