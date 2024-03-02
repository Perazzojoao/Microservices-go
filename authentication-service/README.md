# Authentication Service

Authentication é um microservice especializado em autenticação de usuários. Ao ser requerido, ele pode retornar ou alterar informações do usuário armazenadas em um banco de dados bem como autenticar sua senha.

## Dependências

- `chi`:

    go get github.com/go-chi/chi/v5

- `chi middleware`:

    go get github.com/go-chi/chi/v5/middleware

- `chi cors`:

    go get github.com/go-chi/cors

- `bcrypt`:

    go get golang.org/x/crypto/bcrypt