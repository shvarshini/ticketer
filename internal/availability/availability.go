package availability

import (
	"ticketer/internal/catalog"
	"ticketer/internal/core/lock"
)

type AvailabilityService struct {
	showSeatRepo catalog.ShowSeatRepository
	lockService  lock.LockService
}

func New(showSeatRepo catalog.ShowSeatRepository, lockService lock.LockService) *AvailabilityService {
	if showSeatRepo == nil || lockService == nil {
		panic("Constructor parameter is nil for New AvailabilityService")
	}
	return &AvailabilityService{
		showSeatRepo: showSeatRepo,
		lockService:  lockService,
	}
}

type Service interface {
	GetAvailableSeats(showID string) ([]catalog.ShowSeat, error)
}

func (s *AvailabilityService) GetAvailableSeats(showID string) ([]catalog.ShowSeat, error) {
	return s.showSeatRepo.GetAvailableSeats(showID)
}


