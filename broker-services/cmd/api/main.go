package main

import (
	"fmt"
	"log"
	"net/http"

	"broker/cmd/api/routes"

)

const webPort = "80"

func main() {
	app := routes.Config{}

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
