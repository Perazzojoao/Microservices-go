package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	"authentication/cmd/api/routes"
	"authentication/data"

)

const webPort = "80"

var counts int64

func main() {
	log.Println("Starting authentication service...")

	// Conection to database
	conn := conectToDB()
	if conn == nil {
		log.Fatal("Não foi possível conectar ao banco de dados.")
	}

	// Set up configuration
	app := routes.Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func conectToDB() *sql.DB {
	dns := os.Getenv("DNS")
	for {
		connection, err := openDB(dns)
		if err != nil {
			log.Println("Postgres ainda está indisponível.")
			counts++
		} else {
			log.Println("Postgres conectado com sucesso.")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Tentando reconectar em 2 segundos...")
		time.Sleep(2 * time.Second)
	}
}
