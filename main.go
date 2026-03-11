package main

import (
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/shivsperfect/hotel-reservation/api"
	"github.com/shivsperfect/hotel-reservation/api/middleware"
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
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = &db.Store{
			User:  userStore,
			Hotel: hotelStore,
			Room:  roomStore,
		}
		// handlers initialization
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		authHandler  = api.NewAuthHandler(userStore)
		// fiber app initialization
		app   = fiber.New(config)
		auth  = app.Group("api")
		apiV1 = app.Group("api/v1", middleware.JWTAuthentication)
	)

	// auth handlers
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versioned API routes
	// user handlers
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel handlers
	apiV1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotels/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	if err := app.Listen(*listenAddr); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
