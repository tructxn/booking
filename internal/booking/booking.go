package booking

import (
	"strconv"
	"time"
)

type Entity struct {
	ID           string `json:"id"` // rule base = hotelID checkInDate + checkOutDate + roomID for idempotent
	HotelID      string
	CheckInDate  time.Time `json:"check_in_date"`
	CheckOutDate time.Time `json:"check_out_date"`
	RoomId       int
	Status       BookingStatus `json:"status"`
	Description  string        `json:"description"`
	CreatAt      time.Time     `json:"created_at"`
	UpdateAt     time.Time     `json:"updated_at"`
}

type Dto struct {
	ID           string `json:"id"` // rule base = hotelID checkInDate + checkOutDate + roomID for idempotent
	HotelID      string
	CheckInDate  time.Time
	CheckOutDate time.Time
	RoomId       int
	Status       BookingStatus `json:"status"`
	Description  string        `json:"description"`
}

type BookingStatus string

const (
	StatusInit           BookingStatus = "init"
	StatusProcessPayment BookingStatus = "process_payment"
	StatusSuccessPayment BookingStatus = "success_payment"
	StatusCancel         BookingStatus = "cancel"
	StatusPendingPayment BookingStatus = "pending_payment"
)

func BuildRoomId(hotelID string, checkInDate, checkOutDate time.Time, roomID int) string {
	checkInDateString := GetDateString(checkInDate)
	checkOutDateString := GetDateString(checkOutDate)

	//join all string to make idempotent with split --
	return hotelID + "--" + checkInDateString + "--" + checkOutDateString + "--" + strconv.Itoa(roomID)
}

func GetDateString(date time.Time) string {
	return date.Format("2006-01-02")

}

func DtoToEntity(b *Dto) Entity {
	return Entity{
		ID:           b.ID,
		HotelID:      b.HotelID,
		CheckInDate:  b.CheckInDate,
		CheckOutDate: b.CheckOutDate,
		RoomId:       b.RoomId,
		Status:       b.Status,
		Description:  b.Description,
	}
}

func EntityToDto(b *Entity) Dto {
	return Dto{
		ID:           b.ID,
		HotelID:      b.HotelID,
		CheckInDate:  b.CheckInDate,
		CheckOutDate: b.CheckOutDate,
		RoomId:       b.RoomId,
		Status:       b.Status,
		Description:  b.Description,
	}
}

type AvailableFilter struct {
	HotelID      string
	CheckInDate  time.Time
	CheckOutDate time.Time
}
