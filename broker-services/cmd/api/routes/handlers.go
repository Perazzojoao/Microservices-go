package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

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

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var RequestPayload RequestPayload

	err := app.readJson(w, r, &RequestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch RequestPayload.Action {
	case "auth":
		app.authenticate(w, RequestPayload.Auth)
	case "log":
		app.logItem(w, RequestPayload.Log)
	default:
		app.errorJSON(w, errors.New("invalid action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	// Criar um JSON para enviar ao serviço de log
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	loggerServiceURL := "http://logger-service/log"

	// Enviar o JSON para o serviço de log
	request, err := http.NewRequest("POST", loggerServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// Receber a resposta do serviço de log com o statuscode correto
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling log service"), http.StatusInternalServerError)
		return
	}

	// Criar uma vaeiável para armazenar a resposta do serviço de logger
	var payload = jsonResponse{
		Error:   false,
		Message: "Logged",
	}

	// Escrever a resposta da requisição no formato JSON e definir o status da resposta
	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// Criar um JSON para enviar ao serviço de autenticação
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// Enviar o JSON para o serviço de autenticação
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// Receber a resposta do serviço de autenticação com o statuscode correto
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"), http.StatusUnauthorized)
		return
	}

	// Criar uma vaeiável para armazenar a resposta do serviço de autenticação
	var jasonFromService jsonResponse

	// Decodificar a resposta do serviço de autenticação
	err = json.NewDecoder(response.Body).Decode(&jasonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jasonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload = jsonResponse{
		Error:   false,
		Message: "Authenticated!",
		Data:    jasonFromService.Data,
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}

}
