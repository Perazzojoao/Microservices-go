# Microservises With Go

Projeto desenvolvido durante o curso "Working with Microservices in Go (Golang)" na Udemy.

## Microservices

Neste curso serão desenvolvidos os seguintes microsevices:

- `Brocker`: Ponto de entrada opcional para todos os microservises

- `Authentication`: Autentificação com postgres

- `Logger`: MongoDb

- `Mail`: Envia e-mails com templates específicos

- `Listener`: Recebe mensagens em RabbitMq e inicializa processos

## Ferramentas

Para auxiliar no desenvolvimento serão utilizada as seguintes ferramentas:

- `Docker`: Conteinerizar nossos microservices desenvolvidos

- `GNU Make`: Auxilia na geração de Make files para automatizar compilação de projetos

- `Kubernetes`: Agrupar todos os containers dos microservices em apenas um só container.

## Comunicação

Os microservices desenvolvidos durante o curso se comunicarão entre sí através de:

- `JSON`

- Mandando e recebendo informações utilizando `RPC`

- Mandando e recebendo informações utilizando `gRPC`

- Iniciando e respondendo a eventos através de "Advanced Message Queuing Protocol" (`AMQP`)