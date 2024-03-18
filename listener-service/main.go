package main

import (
	"log"
	"os"

	"listener/database"
	"listener/event"
)

func main() {
	// Conectar com rabbitmq
	rabbitConn, err := database.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	//Começar a ouvir a fila de mensagens
	log.Println("Listening for messages...")

	// Criar um consumidor de eventos
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// Observar fila de mensagens e enviar para o servidor http
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}
