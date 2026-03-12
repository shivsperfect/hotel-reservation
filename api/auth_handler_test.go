package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/shivsperfect/hotel-reservation/types"
)

func insertTestUser(t *testing.T, tdb *testdb) *types.User {
	t.Helper()
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "Shiva",
		LastName:  "Test",
		Email:     "shiva@gmail.com",
		Password:  "password_123",
	})
	if err != nil {
		t.Fatal(err)
	}
	insertedUser, err := tdb.UserStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return insertedUser
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	insertedUser := insertTestUser(t, tdb)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "shiva@gmail.com",
		Password: "password_123",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if authResponse.Token == "" {
		t.Fatal("Expected token to be present in the response")
	}
	if authResponse.User == nil {
		t.Fatal("Expected user to be present in the response")
	}
	if authResponse.User.Email != insertedUser.Email {
		t.Errorf("Expected email %s, got %s", insertedUser.Email, authResponse.User.Email)
	}
	if authResponse.User.FirstName != insertedUser.FirstName {
		t.Errorf("Expected firstName %s, got %s", insertedUser.FirstName, authResponse.User.FirstName)
	}
	if authResponse.User.LastName != insertedUser.LastName {
		t.Errorf("Expected lastName %s, got %s", insertedUser.LastName, authResponse.User.LastName)
	}
}

func TestAuthenticateWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	insertTestUser(t, tdb)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "shiva@gmail.com",
		Password: "wrong_password",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Expected authentication failure status, got %d", resp.StatusCode)
	}

	var errResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	if errResp.Type != "error" {
		t.Fatalf("Expected error response type, got %q", errResp.Type)
	}
	if errResp.Msg != "Invalid Credentials" {
		t.Fatalf("Expected message %q, got %q", "Invalid Credentials", errResp.Msg)
	}
}
