package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/shivsperfect/hotel-reservation/db"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hotelStore db.HotelStore, roomStore db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hotelStore,
		roomStore:  roomStore,
	}
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var queryParams HotelQueryParams
	if err := c.QueryParser(&queryParams); err != nil {
		return err
	}
	fmt.Println(queryParams)
	hotels, err := h.hotelStore.GetAll(c.Context(), nil)
	if err != nil {
		return err
		//return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		//	"error": "failed to fetch hotels",
		//})
	}

	return c.JSON(hotels)
}
