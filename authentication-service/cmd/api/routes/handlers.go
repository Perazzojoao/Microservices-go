package routes

import (
	"bytes"
	"encoding/json"
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

	// Verificar se o usuário existe
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid email or password"))
		return
	}

	// Verificar se a senha está correta
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid email or password"))
		return
	}

	// Envia o registro para o serviço de log
	err = app.LogRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Escreve a resposta
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in as %s", user.Email),
		Data:    user,
	}

	// Envia a resposta
	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) LogRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	loggServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", loggServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
