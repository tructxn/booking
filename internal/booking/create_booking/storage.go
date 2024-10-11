package create_booking

import (
	"demo/booking/internal/booking"
	"github.com/gofiber/fiber/v2/log"
)

type Saver interface {
	SaveBookings(dtos []booking.Dto) ([]booking.Dto, error)
}

type SaverImpl struct {
	Repo booking.Repository // Add the repository dependency here
}

func NewSaver(repo booking.Repository) Saver {
	return &SaverImpl{
		Repo: repo,
	}
}

func (g *SaverImpl) SaveBookings(dtos []booking.Dto) ([]booking.Dto, error) {
	// Implement logic to create booking in the repository

	// for loop dto, convert to list entites and save it.

	entities := make([]booking.Entity, len(dtos))
	for i, dto := range dtos {
		entities[i] = booking.DtoToEntity(&dto)
	}

	savedEntities, err := g.Repo.Save(&entities)
	if err != nil {
		log.Errorf("Error creating booking: %v", err)
		return nil, err
	}
	var bookingDtos []booking.Dto
	for _, entity := range savedEntities {
		bookingDtos = append(bookingDtos, booking.EntityToDto(&entity))
	}

	return bookingDtos, nil
}
