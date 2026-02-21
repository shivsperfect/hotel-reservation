package main

import (
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/shivsperfect/hotel-reservation/api"
	"github.com/shivsperfect/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const dburi = "mongodb://localhost:27017"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(fiber.Map{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "Address to listen on")
	flag.Parse()

	clientOpts := options.Client().ApplyURI(dburi)
	clientOpts.SetBSONOptions(&options.BSONOptions{
		ObjectIDAsHexString: true,
	})
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiV1 := app.Group("api/v1")

	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	if err := app.Listen(*listenAddr); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
