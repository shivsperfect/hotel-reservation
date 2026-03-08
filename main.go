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

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(fiber.Map{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "Address to listen on")
	flag.Parse()

	clientOpts := options.Client().ApplyURI(db.DBURI)
	clientOpts.SetBSONOptions(&options.BSONOptions{
		ObjectIDAsHexString: true,
	})
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	var (
		// handlers initialization
		userHandler  = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		hotelHandler = api.NewHotelHandler(hotelStore, roomStore)

		// fiber app initialization
		app   = fiber.New(config)
		apiV1 = app.Group("api/v1")
	)

	// user handlers
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel handlers
	apiV1.Get("/hotels", hotelHandler.HandleGetHotels)

	if err := app.Listen(*listenAddr); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
