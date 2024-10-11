package api

import (
	"demo/booking/internal/hotel"
	"github.com/gofiber/fiber/v2"
	_ "strconv"
)

type HotelRouter struct {
	HotelStorage hotel.HotelStorage
	RoomStorage  hotel.RoomStorage
}

func NewHotelRouter(hotelStorage hotel.HotelStorage, roomStorage hotel.RoomStorage) *HotelRouter {
	return &HotelRouter{
		HotelStorage: hotelStorage,
		RoomStorage:  roomStorage,
	}
}

// GetHotels godoc
// @Summary Get all hotels
// @Description Get all hotels
// @Tags hotels
// @Accept json
// @Produce json
// @Success 200 {array} hotel.Hotel
// @Router /hotels/ [get]
func (pr *HotelRouter) GetHotels(c *fiber.Ctx) error {
	hotels, err := pr.HotelStorage.AllHotels()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(hotels)
}

// GetHotelById godoc
// @Summary Get hotel by id
// @Description Get hotel by id
// @Tags hotels
// @Accept json
// @Produce json
// @Param id path string true "Hotel ID"
// @Success 200 {object} hotel.Hotel
// @Router /hotels/{id} [get]
func (pr *HotelRouter) GetHotelById(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := pr.HotelStorage.HotelById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(hotel)
}

// GetRoomByHotelId godoc
// @Summary Get room by hotel id
// @Description Get room by hotel id
// @Tags hotels
// @Accept json
// @Produce json
// @Param hotelId path string true "Hotel ID"
// @Success 200 {array} hotel.Room
// @Router /hotels/{hotelId}/rooms [get]
func (pr *HotelRouter) GetRoomByHotelId(c *fiber.Ctx) error {
	hotelId := c.Path("hotelId")
	rooms, err := pr.RoomStorage.RoomsByHotelId(hotelId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(rooms)
}
