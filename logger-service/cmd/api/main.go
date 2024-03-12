package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log-service/cmd/api/routes"
	"log-service/data"

)

var client *mongo.Client

func main() {
	// Connect to MongoDB
	mongoClient, err := conectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// Criando contexto para cancelar a conexão com o banco de dados
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Fecha a conexão com o banco de dados
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := routes.Config{
		Models: data.New(client),
	}

	// Iniciar servidor
	app.Serve()
}

func conectToMongo() (*mongo.Client, error) {
	// Criar opções de conexão
	clientOptions := options.Client().ApplyURI(routes.MongoURL)
	clientOptions.Auth = &options.Credential{
		Username: "admin",
		Password: "password",
	}

	// Conectar ao servidor
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB: ", err)
		return nil, err
	}
	return c, nil
}
