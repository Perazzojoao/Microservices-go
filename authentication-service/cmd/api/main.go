package main

import (
	"log"
	"net/http"

	"authentication/cmd/api/routes"
)

const webPort = "80"

func main() {
	log.Println("Starting authentication service...")

	// Conection to database

	// Set up configuration
	app := routes.Config{}

	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.Routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
