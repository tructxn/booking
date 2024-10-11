package main

import (
	"context"
	_ "demo/booking/docs"
	"demo/booking/internal"
	"demo/booking/internal/storage"
	"demo/booking/pkg/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Payment API
// @version 1.0
// @description This is payment API
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func setUpRoutes(app *fiber.App, pr *api.BookingRoute) {
	app.Get("/booking/:hotelId", pr.GetAllBookings)
	app.Post("/booking/", pr.CreateBooking)
	app.Get("/booking/:bookingId", pr.GetBookingByID)
	app.Post("/booking/available", pr.GetAvailableRoom)
}

func setUpHealthcheck(app *fiber.App, pr *api.HealthcheckRoute) {
	app.Get("/hello", pr.Hello)
}

// @title Payment API
// @version 1.0
// @description This is payment API
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func setUpHotelRoutes(app *fiber.App, pr *api.HotelRouter) {
	app.Get("/hotels/", pr.GetHotels)
}
func setUpDocs(app *fiber.App) {
	// Create routes group.
	route := app.Group("/docs")
	// Routes for GET method:
	route.Get("*", swagger.HandlerDefault) // get one user by ID}
}

func main() {
	defer storage.Close()
	pr, hr, health := internal.Launch()
	app := fiber.New()
	setUpRoutes(app, pr)
	setUpHotelRoutes(app, hr)
	setUpHealthcheck(app, health)
	setUpDocs(app)
	app.Use(cors.New())
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Start the server in a goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	// Listen for the interrupt signal
	shutdown(quit, app)
}

func shutdown(quit chan os.Signal, app *fiber.App) {
	<-quit
	log.Println("Shutting down server...")
	// Set a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}
	log.Println("Server stopped gracefully.")
}
