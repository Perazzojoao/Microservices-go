package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	// queueName string
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setUp()
	if err != nil {
		return &Consumer{}, err
	}

	return &consumer, nil
}

func (consumer *Consumer) setUp() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"logs_topic", // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			payload := Payload{}
			json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()
	fmt.Println(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	// Criar um JSON para enviar ao serviço de log
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	loggerServiceURL := "http://logger-service/log"

	// Enviar o JSON para o serviço de log
	request, err := http.NewRequest("POST", loggerServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Receber a resposta do serviço de log com o statuscode correto
	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
