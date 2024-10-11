package create_booking

import (
	"demo/booking/internal/booking"
	"demo/booking/internal/hotel"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

// HotelConsumer struct
type HotelConsumer struct {
	HotelID     string
	bookedRoom  map[string]*roaring.Bitmap
	roomList    *roaring.Bitmap
	requestChan chan CreateBookingRequest
}

// Create a new hotel consumer
func NewHotelConsumer(hotelID string,
	entities []booking.Entity,
	roomListRaw []hotel.Room,
) *HotelConsumer {
	roomList := roaring.New()
	for _, room := range roomListRaw {
		roomList.Add(uint32(room.RoomID))
	}
	bookedRoomInit := loadBookedBitMap(entities)
	consumer := &HotelConsumer{
		HotelID:     hotelID,
		bookedRoom:  bookedRoomInit,
		roomList:    roomList,
		requestChan: make(chan CreateBookingRequest, 1000),
	}

	go consumer.Run() // Start the consumer thread
	log.Infof("number of booked room %d, number of room ", len(bookedRoomInit), len(roomListRaw))
	log.Infof("start hotel consumer %s", hotelID)
	return consumer
}

func loadBookedBitMap(entities []booking.Entity) map[string]*roaring.Bitmap {
	//loop entity then build bit map map
	bookedMap := make(map[string]*roaring.Bitmap)
	for _, entity := range entities {
		startTime := entity.CheckInDate
		endTime := entity.CheckOutDate
		for i := startTime; i.Before(endTime); i = i.AddDate(0, 0, 1) {
			dateString := booking.GetDateString(i)
			// Get the bitmap for the current date
			bitmap, existed := bookedMap[dateString]
			if !existed {
				bitmap = roaring.New()
				bookedMap[dateString] = bitmap
			}
			// Mark the room as booked
			bitmap.Add(uint32(entity.RoomId))
		}
	}
	return bookedMap
}

// Dedicated thread to handle bookings for each hotel
func (hc *HotelConsumer) Run() {
	for bookingRequest := range hc.requestChan {
		log.Infof("start process booking %s", bookingRequest.HotelID)
		available, err := hc.GetAvailability(bookingRequest)
		// If available, proceed to book and send the result to Snapshot Consumer
		if err == nil && available != nil {
			hc.UpdateInMemoryBitmap(*available)
			bookingRequest.ResponseChan <- CreateBookingResponse{
				Dto: *available,
				Err: err,
			}
			// Send booking result to Snapshot Consumer
		} else {
			log.Errorf("Error checking availability for hotel %s: %v", hc.HotelID, err)
			fmt.Printf("Booking failed for hotel %s: not enough rooms.\n", hc.HotelID)
			bookingRequest.ResponseChan <- CreateBookingResponse{
				Dto: booking.Dto{},
				Err: fmt.Errorf("Booking failed for hotel %s: not enough rooms", hc.HotelID),
			}
		}
	}
}

// Update in-memory bitmap (mark rooms as booked)
func (hc *HotelConsumer) UpdateInMemoryBitmap(dto booking.Dto) {
	startTime := dto.CheckInDate
	endTime := dto.CheckOutDate
	for i := startTime; i.Before(endTime); i = i.AddDate(0, 0, 1) {
		dateString := booking.GetDateString(i)
		// Get the bitmap for the current date
		bitmap, existed := hc.bookedRoom[dateString]
		if !existed {
			bitmap = roaring.New()
			hc.bookedRoom[dateString] = bitmap
		}
		// Mark the room as booked
		bitmap.Add(uint32(dto.RoomId))
	}
}

func (hc *HotelConsumer) GetAvailabilities(request booking.AvailableFilter) ([]string, error) {
	availabilities, err := hc.getAvailabilities(request.CheckInDate, request.CheckOutDate)
	if err != nil || availabilities.IsEmpty() {
		return nil, err
	}

	return make([]string, availabilities.GetCardinality()), nil
}
func (hc *HotelConsumer) getAvailabilities(startTime, endTime time.Time) (*roaring.Bitmap, error) {
	// Check if the hotel has enough rooms available
	// Check the in-memory bitmap for available rooms
	// If enough rooms are available, return true
	// Otherwise, return false

	availableRoom := hc.roomList.Clone()
	for i := startTime; i.Before(endTime); i = i.AddDate(0, 0, 1) {
		dateString := booking.GetDateString(i)
		// Get the bitmap for the current date
		bitmap, existed := hc.bookedRoom[dateString]
		if !existed {
			bitmap = roaring.New()
			hc.bookedRoom[dateString] = bitmap
		}
		// get diff between all room and booked room to findout available room
		availableRoom.AndNot(bitmap)
		if availableRoom.IsEmpty() {
			return nil, fmt.Errorf("No rooms available on %s", dateString)
		}
	}

	return availableRoom, nil
}

func (hc *HotelConsumer) GetAvailability(request CreateBookingRequest) (*booking.Dto, error) {
	availableRoom, err := hc.getAvailabilities(request.CheckInDate, request.CheckOutDate)
	if err != nil || availableRoom.IsEmpty() {
		return nil, err
	}

	roomId := availableRoom.Minimum()
	return &booking.Dto{
		ID:           booking.BuildRoomId(hc.HotelID, request.CheckInDate, request.CheckOutDate, int(roomId)),
		HotelID:      hc.HotelID,
		CheckInDate:  request.CheckInDate,
		CheckOutDate: request.CheckOutDate,
		RoomId:       int(roomId),
		Status:       booking.StatusInit,
		Description:  request.Description,
	}, nil
}
