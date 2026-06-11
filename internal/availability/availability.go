package availability

import (
	"ticketer/internal/catalog"
)

type AvailabilityService struct {
	showSeatRepo catalog.ShowSeatRepository
}

func New(showSeatRepo catalog.ShowSeatRepository) *AvailabilityService {
	if showSeatRepo == nil {
		panic("Constructor parameter is nil for New AvailabilityService")
	}
	return &AvailabilityService{
		showSeatRepo: showSeatRepo,
	}
}

type Service interface {
	GetAvailableSeats(showID string) ([]catalog.ShowSeat, error)
}

func (s *AvailabilityService) GetAvailableSeats(showID string) ([]catalog.ShowSeat, error) {
	return s.showSeatRepo.GetAvailableSeats(showID)
}


