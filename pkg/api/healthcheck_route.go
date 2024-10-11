package api

import (
	"demo/booking/internal/booking"
	"demo/booking/internal/booking/create_booking"
	"github.com/gofiber/fiber/v2"
	"log"
	_ "strconv"
)

type HealthcheckRoute struct {
	bookingFetcher booking.Fetcher
}

func NewHealthcheckRoute() *HealthcheckRoute {
	return &HealthcheckRoute{}
}

var msgChannel = make(chan create_booking.CreateBookingRequest, 10000)

// Initialize the consumer goroutine to process messages from the channel
func init() {
	go messageConsumer()
}

// Consumer function that processes the messages from the channel
func messageConsumer() {
	for msg := range msgChannel {
		// Here you can process the message, log it, or do any other operation
		log.Printf("Message received: %s", msg)
	}
}

// @title helloWorld
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @host localhost:3000
// @path /hello
// @BasePath /
func (pr *HealthcheckRoute) Hello(c *fiber.Ctx) error {
	//return c.SendString("fiber")
	msgChannel <- create_booking.CreateBookingRequest{}
	// Respond to the client (optional, just to indicate the message was sent)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Message sent to the channel",
	})
}
