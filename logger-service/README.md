# Logger service

Serviço responsável por registrar todas as ocorrências entre os microservices e armazená-las em um banco de dados. 

## MongoDB

MongoDB é o banco de dados em que o logger-service irá se conectar e armazenar seus dados.

## Dependências

- `chi`:

    go get github.com/go-chi/chi/v5

- `chi middleware`:

    go get github.com/go-chi/chi/middleware

- `chi cors`:

    go get github.com/go-chi/cors

- `mongo-driver`:

    go get go.mongodb.org/mongo-driver/mongo

- `grpc`: Necessário para servidores gRpc

    go get google.golang.org/grpc