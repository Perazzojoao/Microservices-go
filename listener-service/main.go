package main

import (
	"log"
	"os"

	"listener/database"
)

func main() {
	// Conectar com rabbitmq
	rabbitConn, err := database.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	//Começar a ouvir a fila de mensagens

	// Criar um servidor http

	// Observar fila de mensagens e enviar para o servidor http
}
