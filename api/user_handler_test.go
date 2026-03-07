package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/shivsperfect/hotel-reservation/db"
	"github.com/shivsperfect/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testDbUri = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) tearDown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	clientOpts := options.Client().ApplyURI(testDbUri)
	clientOpts.SetBSONOptions(&options.BSONOptions{
		ObjectIDAsHexString: true,
	})
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatal("Failed to connect to test database: ", err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}

}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)
	params := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "some@gmailco.com",
		Password:  "password123",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Errorf("Error making request: %v", err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("Expected user ID to be set, got empty string")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Expected encrypted password to be empty in response")
	}
	if user.Email != params.Email {
		t.Errorf("Expected Firstname %s got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("Expected lastname %s got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("Expected email %s got %s", params.Email, user.Email)
	}
}
