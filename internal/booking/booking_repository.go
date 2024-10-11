package booking

import (
	"demo/booking/internal/storage"
	"time"
)

type Repository interface {
	Create(booking *Entity) (Entity, error)
	Save(bookings *[]Entity) ([]Entity, error)
	FindAll() ([]Entity, error)
	FindByID(uid uint) (Entity, error)
	UpdateStatus(id uint, status string) error
	FindFutureBookedByHotelId(hotelID string) ([]Entity, error)
}

func NewRepository(db *storage.DatabaseInstance) Repository {
	return &RepositoryImpl{
		db: db, // Initialize the db field with the provided database instance
	}
}

type RepositoryImpl struct {
	db *storage.DatabaseInstance
}

func (pr *RepositoryImpl) Create(booking *Entity) (Entity, error) {
	err := pr.db.Db.Create(booking).Error
	return *booking, err
}

func (pr *RepositoryImpl) Save(bookings *[]Entity) ([]Entity, error) {
	if err := pr.db.Db.Save(bookings).Error; err != nil {
		return nil, err // Return nil for the bookings and the error
	}
	// If successful, return the saved bookings (now with IDs populated)
	return *bookings, nil // Dereference the pointer to return the slice
}

func (pr *RepositoryImpl) FindAll() ([]Entity, error) {
	var Bookings []Entity
	err := pr.db.Db.Find(&Bookings).Error
	return Bookings, err
}

func (pr *RepositoryImpl) FindByID(id uint) (Entity, error) {
	var Booking Entity
	err := pr.db.Db.Where("id = ?", id).Find(&Booking).Error
	return Booking, err
}

func (pr *RepositoryImpl) UpdateStatus(id uint, status string) error {
	return pr.db.Db.Model(&Entity{}).Where("id = ?", id).Update("status", status).Error
}

func (pr *RepositoryImpl) FindFutureBookedByHotelId(hotelID string) ([]Entity, error) {
	var Bookings []Entity
	err := pr.db.Db.Where("hotel_id = ? and status != ? and check_out_date > ? ", hotelID, StatusCancel, time.Now()).Find(&Bookings).Error
	return Bookings, err
}
