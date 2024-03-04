package routes

import (
	"errors"
	"fmt"
	"net/http"

)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid email or password"))
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid email or password"))
		return
	}

	payload := jsonResponse{
		Error:  false,
		Message: fmt.Sprintf("Logged in as %s", user.Email),
		Data: 	user,
	}

	app.writeJSON(w, http.StatusOK, payload)
}
