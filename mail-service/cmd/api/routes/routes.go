package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Config struct{}

const WebPort = "80"

func (app *Config) Routes() http.Handler {
	mux := chi.NewRouter()

	// Especificar permissões de CORS
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	mux.Use(middleware.Heartbeat("/ping"))

	return mux
}

func (app *Config) Serve() {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WebPort),
		Handler: app.Routes(),
	}

	log.Println("Starting server on port", WebPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Server error: ", err)
	}
}
