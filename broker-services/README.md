# Broker Service

Brocker é um microservice especializado em receber requisições do frontend e redirecioná-las para um microservice diretamente responsável.

## Dependências

- `chi`:

    go get github.com/go-chi/chi/v5

- `chi middleware`:

    go get github.com/go-chi/chi/v5/middleware

- `chi cors`:

    go get github.com/go-chi/cors

- `rabbitmq`: plugin para db rabbitmq

    go get github.com/rabbitmq/amqp091-go

- `grpc`: Necessário para servidores gRpc

    go get google.golang.org/grpc