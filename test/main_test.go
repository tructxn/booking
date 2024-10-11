package test

import (
	"bytes"
	"context"
	"demo/booking/internal"
	"demo/booking/internal/booking/create_booking"
	"demo/booking/internal/hotel"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	dbTest "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	//Catching all panics to once again make sure that shutDown is successfully run
	defer func() {
		if r := recover(); r != nil {
			shutDown()
			fmt.Println("Panic")
		}
	}()
	SetupTestContainersAndServer()
	code := m.Run()
	shutDown()
	os.Exit(code)
}

func shutDown() {
	// Shut down containers
}

func SetupTestContainersAndServer() func() {
	// Start the PostgreSql container
	ctx := context.Background()

	dbName := "testdb"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := dbTest.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		dbTest.WithDatabase(dbName),
		dbTest.WithUsername(dbUser),
		dbTest.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatal("failed to start container: %s", err)
	}

	// Get the PostgreSQL port
	port, _ := postgresContainer.MappedPort(ctx, "5432")
	println(port)

	parts := strings.Split(string(port), "/")

	// Extract the numeric part
	portStr := parts[0]

	// Convert the numeric part to an integer
	port_number, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("Error parsing port:", err)
		panic("invalid port number")
	}
	//fmt.Printf("conn : %s : port : %d     ", connStr, port)
	// Connect to the PostgreSQL database
	dsn := fmt.Sprintf("host=%s user=user password=password dbname=testdb port=%d sslmode=disable TimeZone=Asia/Shanghai", "localhost", port_number)
	os.Setenv("DSN", dsn)
	// Wait for the Fiber app server to start
	time.Sleep(2 * time.Second)
	pr, _, hr := internal.Launch()

	app := fiber.New()
	app.Get("/booking/:hotelId", pr.GetAllBookings)
	app.Post("/booking/", pr.CreateBooking)
	app.Get("/hello", hr.Hello)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Start the server in a goroutine
	go func() {
		if err := app.Listen(":3009"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Clean up the container
	return func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}

	// Listen for the interrupt signal
}

func TestBooking(t *testing.T) {
	hotels := hotel.GenHotel(internal.NumberOfHotels)
	today := time.Now()
	//convert request to data
	t.Run("InitBooking", func(t *testing.T) {
		for _, h := range hotels {
			request := create_booking.CreateBookingRequest{
				HotelID:      h.ID,
				CheckInDate:  today,
				CheckOutDate: today.AddDate(0, 0, 1),
				RoomTypeId:   10,
				Description:  "",
			}
			data, _ := json.Marshal(request)
			url := "http://localhost:3009/booking/"
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
			if err != nil {
				log.Fatal("Post request failed:", err)
			}
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			//
			body, errRes := io.ReadAll(resp.Body)
			assert.NoError(t, errRes)
			created := convert(body)
			assert.Equal(t, created.HotelID, request.HotelID)
			assert.Equal(t, created.RoomId, 0)
		}
	})

	t.Run("InitBookingRound2", func(t *testing.T) {
		for _, h := range hotels {
			request := create_booking.CreateBookingRequest{
				HotelID:      h.ID,
				CheckInDate:  today,
				CheckOutDate: today.AddDate(0, 0, 1),
				RoomTypeId:   10,
				Description:  "",
			}
			data, _ := json.Marshal(request)
			url := "http://localhost:3009/booking/"
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
			if err != nil {
				log.Fatal("Post request failed:", err)
			}
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			//
			body, errRes := io.ReadAll(resp.Body)
			assert.NoError(t, errRes)
			created := convert(body)
			assert.Equal(t, created.HotelID, request.HotelID)
			assert.Equal(t, created.RoomId, 1)
		}
	})

	h := hotels[0]
	request := create_booking.CreateBookingRequest{
		HotelID:      h.ID,
		CheckInDate:  today,
		CheckOutDate: today.AddDate(0, 0, 1),
		RoomTypeId:   10,
		Description:  "",
	}

	data, _ := json.Marshal(request)
	t.Run("InitBookingRound3", func(t *testing.T) {
		for i := 2; i < internal.NumberOfRooms; i++ {
			url := "http://localhost:3009/booking/"
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
			if err != nil {
				log.Fatal("Post request failed:", err)
			}
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			body, errRes := io.ReadAll(resp.Body)
			assert.NoError(t, errRes)
			created := convert(body)
			assert.Equal(t, created.HotelID, request.HotelID)
			assert.Equal(t, created.RoomId, i)
		}
	})

	t.Run("InitBookingRound4", func(t *testing.T) {
		for i := internal.NumberOfRooms; i < internal.NumberOfRooms*2; i++ {
			url := "http://localhost:3009/booking/"
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
			if err != nil {
				log.Fatal("Post request failed:", err)
			}
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		}
		//
	})
}

func convert(body []byte) create_booking.CreateBookingResponse {
	var dto create_booking.CreateBookingResponse
	// Unmarshal the JSON data into the struct
	err := json.Unmarshal(body, &dto)
	if err != nil {
		log.Fatal(err)
	}
	return dto
}
