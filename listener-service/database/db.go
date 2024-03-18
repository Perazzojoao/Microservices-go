package database

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// Não continuar até que a conexão seja estabelecida
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("Unable to connect to RabbitMQ after 5 retries. Exiting...")
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Failed to connect to RabbitMQ. Retrying in", backOff, "seconds")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
