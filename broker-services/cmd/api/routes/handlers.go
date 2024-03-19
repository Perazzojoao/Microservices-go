package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"broker/event"
	"broker/logs"

)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
		app.logItemViaRPC(w, RequestPayload.Log)
	case "mail":
		app.sendMail(w, RequestPayload.Mail)
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
		app.errorJSON(w, errors.New("error calling log service"))
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

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	// Criar um JSON para enviar ao serviço de email
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// Enviar o JSON para o serviço de email
	mailServiceURL := "http://mailer-service/send"
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
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

	// Receber a resposta do serviço de email com o statuscode correto
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	// Criar uma variável para armazenar a resposta do serviço de email
	var payload = jsonResponse{
		Error:   false,
		Message: "Mail sent to " + msg.To,
	}

	// Escrever a resposta da requisição no formato JSON e definir o status da resposta
	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToTheQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "logged via RabbitMQ",
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}
}

func (app *Config) pushToTheQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.Marshal(&payload)
	err = emitter.Emit(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) logViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if requestPayload.Action != "log" {
		app.errorJSON(w, errors.New("invalid action"))
		return
	}

	// Set up a connection to the gRpc server.
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	// Create a new gRpc client
	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call the gRpc server WriteLog method
	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Create a response payload and write it to the response
	var payload = jsonResponse{
		Error:   false,
		Message: "Logged via gRPC",
	}
	app.writeJSON(w, http.StatusOK, payload)
}
