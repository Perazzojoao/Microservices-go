package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"log-service/data"

)

const (
	WebPort  = "80"
	RpcPort  = "5001"
	MongoURL = "mongodb://mongo:27017"
	GRpcPort = "50001"
)

type Config struct {
	Models data.Models
}

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
	mux.Post("/log", app.WriteLog)
	return mux
}

func (app *Config) Serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WebPort),
		Handler: app.Routes(),
	}

	fmt.Printf("Server is running on port %s\n", WebPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}