package main

import (
	"mailer-service/cmd/api/routes"
)

func main() {
	app := routes.Config{}

	app.Serve()
}
