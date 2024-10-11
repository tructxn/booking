package create_booking

import (
	"demo/booking/internal/booking"
	"demo/booking/internal/hotel"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

type BookingCreator interface {
	ProcessBooking(CreateBookingRequest) CreateBookingResponse
	FindAvailableSlot(filter *booking.AvailableFilter) ([]string, error)
}

type CreateBookingRequest struct {
	HotelID      string
	CheckInDate  time.Time
	CheckOutDate time.Time
	RoomTypeId   int
	Description  string                     `json:"description"`
	ResponseChan chan CreateBookingResponse `json:"-"`
}

type CreateBookingResponse struct {
	booking.Dto
	Err error
}

type BookingCreatorImpl struct {
	hotelConsumers map[string]*HotelConsumer
	storageSaver   Saver
	hotelFetcher   hotel.HotelStorage
	bookingRepo    booking.Repository
	channelMap     map[string]chan CreateBookingRequest
}

func NewBookingCreator(storageSaver Saver,
	hotelFetcher hotel.HotelStorage,
	bookingRepo booking.Repository,
	roomRepo hotel.RoomStorage) *BookingCreatorImpl {
	hotels, err := hotelFetcher.AllHotels()
	if err != nil {
		log.Fatal("Init hotel error ")
	}
	hotelConsumers := make(map[string]*HotelConsumer, len(hotels))
	channelMap := make(map[string]chan CreateBookingRequest)

	for _, info := range hotels {
		booked, err := bookingRepo.FindFutureBookedByHotelId(info.ID)
		if err != nil {
			log.Fatalf("cannot load booked slot ")
		}
		roomList, err := roomRepo.RoomsByHotelId(info.ID)
		consumer := NewHotelConsumer(info.ID, booked, roomList)
		hotelConsumers[info.ID] = consumer
		channelMap[info.ID] = consumer.requestChan
	}
	return &BookingCreatorImpl{
		hotelConsumers: hotelConsumers,
		storageSaver:   storageSaver,
		hotelFetcher:   hotelFetcher,
		bookingRepo:    bookingRepo,
		channelMap:     channelMap,
	}
}

func (c *BookingCreatorImpl) ProcessBooking(createRequest CreateBookingRequest) CreateBookingResponse {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("error request %s", r)
		}
	}()

	// Ensure the channel exists
	channel, exists := c.channelMap[createRequest.HotelID]
	if !exists {
		return CreateBookingResponse{
			Err: fmt.Errorf("Hotel %s not found", createRequest.HotelID),
		}
	}

	responseChan := make(chan CreateBookingResponse, 1) // Buffered channel with size 1
	// Attach the response channel to the request
	createRequest.ResponseChan = responseChan

	go func() {
		select {
		case channel <- createRequest:
			log.Infof("publish success request %s", createRequest.HotelID)
		default:
			log.Errorf(" error request, ")
		}
	}()

	// Wait for the response
	response := <-responseChan
	return response
	// Send booking to internal queue for async processin
}

func (c *BookingCreatorImpl) FindAvailableSlot(filter *booking.AvailableFilter) ([]string, error) {
	// Find the consumer for the hotel
	hotelConsumer, ok := c.hotelConsumers[filter.HotelID]
	if !ok {
		return nil, fmt.Errorf("Hotel %s not found", filter.HotelID)
	}
	return hotelConsumer.GetAvailabilities(*filter)
}
