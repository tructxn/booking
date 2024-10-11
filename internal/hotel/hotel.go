package hotel

import "time"

type Hotel struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Location    string
	Description string
	UpdateDate  time.Time
	CreateDate  time.Time
}

type Room struct {
	HotelID     string `gorm:"primaryKey;column:hotel_id"` // Composite key part 1
	RoomID      int    `gorm:"primaryKey;column:room_id"`  // Composite key part 2
	RoomType    int
	Price       float64
	Description string
	UpdateDate  time.Time
	CreateDate  time.Time
}
