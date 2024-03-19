package routes

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"

	"log-service/data"
	"log-service/logs"

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

func (app *Config) RpcListen() error {
	log.Println("Starting RPC server on port ", RpcPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", RpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func (app *Config) GrpcListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", GRpcPort))
	if err != nil {
		log.Fatalf("failed to listen for gRpc: %v", err)
	}
	defer listen.Close()

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
	log.Println("gRPC server started on port ", GRpcPort)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRpc: %v", err)
	}

}
