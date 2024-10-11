package hotel

import "demo/booking/internal/storage"

type HotelStorage interface {
	AllHotels() ([]Hotel, error)
	HotelById(id string) (Hotel, error)
	saves(hotels []Hotel) error
}

func NewHotelStorage(db *storage.DatabaseInstance) HotelStorage {
	return &HotelStorageImpl{db: db}
}

type HotelStorageImpl struct {
	db *storage.DatabaseInstance
}

func (hs *HotelStorageImpl) AllHotels() ([]Hotel, error) {
	var hotels []Hotel
	err := hs.db.Db.Find(&hotels).Error
	return hotels, err
}

func (hs *HotelStorageImpl) HotelById(id string) (Hotel, error) {
	var hotel Hotel
	err := hs.db.Db.Where("id = ?", id).Find(&hotel).Error
	return hotel, err
}

func (hs *HotelStorageImpl) saves(hotels []Hotel) error {
	err := hs.db.Db.Save(hotels).Error
	return err
}

type RoomStorage interface {
	RoomsByHotelId(hotelId string) ([]Room, error)
	RoomById(hotelId string, roomId int) (Room, error)
	Saves(room []Room) error
}

func NewRoomStorage(db *storage.DatabaseInstance) RoomStorage {
	return &RoomStorageImpl{db: db}
}

type RoomStorageImpl struct {
	db *storage.DatabaseInstance
}

func (rs *RoomStorageImpl) Saves(rooms []Room) error {
	err := rs.db.Db.Save(rooms).Error
	return err
}

func (rs *RoomStorageImpl) RoomsByHotelId(hotelId string) ([]Room, error) {
	var rooms []Room
	err := rs.db.Db.Where("hotel_id = ?", hotelId).Find(&rooms).Error
	return rooms, err
}

func (rs *RoomStorageImpl) RoomById(hotelId string, roomId int) (Room, error) {
	var room Room
	err := rs.db.Db.Where("hotel_id = ? and id = ?", hotelId, roomId).Find(&room).Error
	return room, err
}
