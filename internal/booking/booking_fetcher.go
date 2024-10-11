package booking

type Fetcher interface {
	AllBookings() ([]Dto, error)
	BookingById(id uint) (Dto, error)
	FindFutureBookedByHotelId(hotelID string) ([]Dto, error)
}

func NewFetcher(repository Repository) Fetcher {
	return &FetcherImpl{
		Repo: repository,
	}
}

type FetcherImpl struct {
	Repo Repository // Add the repository dependency here
}

func (g *FetcherImpl) FindFutureBookedByHotelId(hotelID string) ([]Dto, error) {
	ids, err := g.Repo.FindFutureBookedByHotelId(hotelID)

	if err != nil {
		return nil, err
	}

	// Map booking Entity to booking DTO
	bookingDtos := make([]Dto, len(ids))
	for i, id := range ids {
		bookingDtos[i] = EntityToDto(&id)
	}

	return bookingDtos, err
}

func (g *FetcherImpl) AllBookings() ([]Dto, error) {
	bookings, err := g.Repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Map booking Entity to booking DTO
	bookingDtos := make([]Dto, len(bookings))
	for i, booking := range bookings {
		bookingDtos[i] = Dto{
			ID: booking.ID,
		}
	}
	return bookingDtos, nil
}

func (g *FetcherImpl) BookingById(id uint) (Dto, error) {
	// Implement logic to retrieve payment by ID from the repository
	booking, err := g.Repo.FindByID(id)
	if err != nil {
		return Dto{}, err
	}

	// Map booking Entity to booking DTO
	bookingDto := Dto{
		ID: booking.ID,
	}
	return bookingDto, nil
}
