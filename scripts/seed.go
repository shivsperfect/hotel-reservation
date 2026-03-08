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

	if err := client.Database(db.DBNAME).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}

func main() {
	seeHotel("Bellucia", "France", 5)
	seeHotel("The cozy hotel", "Netherland", 4)
	seeHotel("Don't die in you sleep", "England", 3)
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
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 149.9,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 199.9,
		},
	}
	insertedHotel, err := hotelStore.Insert(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(context.Background(), &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}
