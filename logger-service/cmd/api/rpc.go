package main

import (
	"context"
	"log"
	"time"

	"log-service/data"
)

// RPCServer é a estrutura que implementa os métodos do servidor RPC
type RPCServer struct{}

// RPCPayload é a estrutura que representa os dados enviado via RPC
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo é o método que escreve o payload para o mongoDB
func (r *RPCServer) LogInfo(payload RPCPayload, resq *string) error {
	colecton := client.Database("logs").Collection("logs")
	_, err := colecton.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log int to mongo: ", err)
		return err
	}

	*resq = "Processed payload via RPC: " + payload.Name
	return nil
}
