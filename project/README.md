# Project

Diretório responsável por conter arquivos de configuração do projeto. É usado, principalmente, como diretório raíz para comandos `docker-compose` e `make`.

## Docker-compose

Arquivo responsável por definir todas as imagens doker do projeto e iniciá-las em conjunto

## Makefile

Arquivo que contém diversos comandos que servem como atalhos e automações para certas tarefas
realizadas com frequência. Contém os seguintes atalhos:

- `up`: Inicia todos os containers

- `down`:Para todos os containers

- `up_build`: Para todos os containers, recompila e inicia todos novamente.

- `build_broker`: Compila projeto broker para o sistema linux

- `start`: Inicia Front end

- `stop`: Para Front end

- `build_front`: Compila Front end