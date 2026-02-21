package main

import (
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/shivsperfect/hotel-reservation/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "Address to listen on")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("api/v1")

	apiV1.Get("/user", api.HandleGetUsers)
	apiV1.Get("/user/:id", api.HandleGetUser)
	err := app.Listen(*listenAddr)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
