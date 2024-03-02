package routes

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hello, World!",
	}

	err := app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}
