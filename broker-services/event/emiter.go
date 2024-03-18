package event

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return declareExchange(channel)
}

func (e *Emitter) Emit(event, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	log.Printf("Publishing to channel")

	err = channel.PublishWithContext(
		context.TODO(), // context
		"logs_topic",    // exchange
		severity,       // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) (*Emitter, error) {
	emmiter := &Emitter{
		connection: conn,
	}

	err := emmiter.setup()
	if err != nil {
		return &Emitter{}, err
	}

	return emmiter, nil
}
