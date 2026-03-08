package db

const (
	DBNAME      = "hotel-reservation"
	DBURI       = "mongodb://localhost:27017"
	TEST_DBNAME = "hotel-reservation-test"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
