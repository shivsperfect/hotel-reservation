package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shivsperfect/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		ID:        "",
		FirstName: "James",
		LastName:  "Bond",
	}
	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(id)
}
