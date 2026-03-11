package main

import (
	"context"
	"log"

	"github.com/shivsperfect/hotel-reservation/db"
	"github.com/shivsperfect/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	userStore  db.UserStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func init() {
	var err error
	clientOpts := options.Client().ApplyURI(db.DBURI)
	client, err = mongo.Connect(clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}

func main() {
	seeHotel("Bellucia", "France", 5)
	seeHotel("The cozy hotel", "Netherland", 4)
	seeHotel("Don't die in you sleep", "England", 3)
	seedUser("shiva", "kumar", "shiva@gmail.com")
}

func seeHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []bson.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size:  "single",
			Price: 99.9,
		},
		{
			Size:  "double",
			Price: 149.9,
		},
		{
			Size:  "deluxe",
			Price: 199.9,
		},
	}
	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func seedUser(fname, lname, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  "password",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}
