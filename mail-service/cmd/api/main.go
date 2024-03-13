package main

import (
	"mailer-service/cmd/api/mailer"
	"mailer-service/cmd/api/routes"

)

func main() {
	app := routes.Config{
		Mailer: mailer.CreateMail(),
	}

	app.Serve()
}
