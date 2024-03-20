package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"broker/cmd/api/routes"
	"broker/database"
)

const webPort = "8080"

func main() {
	// Conectar com rabbitmq
	rabbitConn, err := database.Connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := routes.Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Server is running on port:%s\n", webPort)

	// Define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	// Start server
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
