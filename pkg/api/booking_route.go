package api

import (
	"demo/booking/internal/booking"
	"demo/booking/internal/booking/create_booking"
	"github.com/gofiber/fiber/v2"
	"strconv"
	_ "strconv"
)

type BookingRoute struct {
	bookingFetcher booking.Fetcher
	bookingCreator create_booking.BookingCreator
}

// Initialize the consumer goroutine to process messages from the channel
func init() {
	go messageConsumer()
}

func NewBookingRoute(bookingFetcher booking.Fetcher,
	bookingCreator create_booking.BookingCreator,
) *BookingRoute {

	return &BookingRoute{
		bookingFetcher: bookingFetcher,
		bookingCreator: bookingCreator,
	}
}

// CreateBooking
// @Summary create a booking
// @description create a booking
// @Tags booking
// @Accept json
// @Produce json
// @Param booking body create_booking.CreateBookingRequest true "Booking information"
// @Success 200 {object} booking.Dto
// @Router /booking/ [post]
func (pr *BookingRoute) CreateBooking(c *fiber.Ctx) error {
	//log.Info("Create booking, request received, processing...")
	book := new(create_booking.CreateBookingRequest)
	//// Attempt to parse the request body into the BookingDto struct
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}
	processBooking := pr.bookingCreator.ProcessBooking(*book)
	if processBooking.Err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(processBooking.Err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(processBooking)
}

// GetAvailableRoom
// @Summary Get available room
// @description Get available room
// @Tags booking
// @Accept json
// @Produce json
// @Param filter body booking.AvailableFilter true "Filter information"
// @Success 200 {object} []string
// @Router /booking/available [post]
func (pr *BookingRoute) GetAvailableRoom(c *fiber.Ctx) error {
	availableFilter := new(booking.AvailableFilter)
	if err := c.BodyParser(availableFilter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}
	availableRooms, err := pr.bookingCreator.FindAvailableSlot(availableFilter)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(availableRooms)
}

// GetAllBookings
// @Summary Get all bookings
// @description Get all bookings
// @Tags booking
// @Produce json
// @Success 200 {object} booking.Dto
// @Router /booking/:hotelId [get]
func (pr *BookingRoute) GetAllBookings(c *fiber.Ctx) error {
	hotelId := c.Params("hotelId")
	bookings, err := pr.bookingFetcher.FindFutureBookedByHotelId(hotelId)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(bookings)
}

// GetBookingByID
// @Summary Get booking by id
// @description Get booking by id
// @Tags booking
// @Produce json
// @Param bookingId path string true "Booking ID" // This line is added
// @Success 200 {object} booking.Dto
// @Router /booking/:bookingId [get]
func (pr *BookingRoute) GetBookingByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("bookingId")) // Get the ID from the URL parameter
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	bookingById, err := pr.bookingFetcher.BookingById(uint(id))
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(bookingById)
}
