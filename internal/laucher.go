package internal

import (
	"demo/booking/internal/booking"
	"demo/booking/internal/booking/create_booking"
	"demo/booking/internal/hotel"
	"demo/booking/internal/storage"
	"demo/booking/pkg/api"
	"log"
	"runtime"
)

var NumberOfRooms = 200
var NumberOfHotels = 16

func Launch() (*api.BookingRoute, *api.HotelRouter, *api.HealthcheckRoute) {
	runtime.GOMAXPROCS(8)
	storage.ConnectDb()
	repo := booking.NewRepository(&storage.DB)
	fetcher := booking.NewFetcher(repo)
	hotelStorage := hotel.NewHotelStorage(&storage.DB)
	bookingStorage := create_booking.NewSaver(repo)
	roomStorage := hotel.NewRoomStorage(&storage.DB)

	///storageSaver Saver, hotelFetcher hotel.HotelStorage, bookingRepo booking.Repository
	err := storage.DB.Db.AutoMigrate(&booking.Entity{}, &hotel.Hotel{}, &hotel.Room{})
	if err != nil {
		log.Fatal(err)
		return nil, nil, nil
	}

	//init test data
	initializer := hotel.NewHotelInitializer(hotelStorage, roomStorage)
	err = initializer.InitializeHotel(NumberOfHotels, NumberOfRooms)
	if err != nil {
		log.Fatal(err)
		return nil, nil, nil
	}
	rommRepo := hotel.NewRoomStorage(&storage.DB)
	creator := create_booking.NewBookingCreator(
		bookingStorage,
		hotelStorage,
		repo,
		rommRepo)
	pr := api.NewBookingRoute(fetcher, creator) // Inject concrete implementation
	hr := api.NewHotelRouter(hotelStorage, roomStorage)
	healthR := api.NewHealthcheckRoute()
	return pr, hr, healthR
}
