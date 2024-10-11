package hotel

import (
	"github.com/gofiber/fiber/v2/log"
	"math/rand"
	"strconv"
	"time"
)

type InitializeHotel interface {
	InitializeHotel() error
}

type InitializeHotelImpl struct {
	HotelStorage HotelStorage
	RoomStorage  RoomStorage
}

func NewHotelInitializer(hotelStorage HotelStorage, roomStorage RoomStorage) *InitializeHotelImpl {
	return &InitializeHotelImpl{
		HotelStorage: hotelStorage,
		RoomStorage:  roomStorage,
	}
}

func (hi *InitializeHotelImpl) InitializeHotel(hotelSize int, roomSize int) error {

	hotels := GenHotel(hotelSize)
	_ = hi.HotelStorage.saves(hotels)
	for _, hotel := range hotels {
		rooms := GenRoom(hotel, roomSize)
		_ = hi.RoomStorage.Saves(rooms)
		log.Info("save hotel room : %s %s", hotel.ID, len(rooms))
	}
	return nil
}

func GenRoom(hotel Hotel, count int) []Room {
	rooms := make([]Room, count)

	roomType := rand.Intn(10)
	for j := 0; j < count; j++ {
		room := Room{
			HotelID:  hotel.ID,
			RoomID:   j,
			RoomType: roomType,
		}
		rooms[j] = room
	}
	return rooms
}

func GenHotel(count int) []Hotel {
	hotels := make([]Hotel, count)
	for i := 0; i < count; i++ {
		//convert i to string
		index := "hotel_" + strconv.Itoa(i) // Convert integer to string using strconv.Itoa
		hotel := Hotel{
			ID:          "ID_" + index,
			Location:    "hanoi",
			CreateDate:  time.Now(),
			UpdateDate:  time.Now(),
			Description: "auto generate " + index,
			Name:        "Hotel " + index,
		}

		hotels[i] = hotel
	}
	return hotels
}
