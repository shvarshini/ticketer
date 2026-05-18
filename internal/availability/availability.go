package availability

import (
	"fmt"
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
	LockSeats(showID string, seatIDs []string) error
	ReleaseSeats(showID string, seatIDs []string) error
	BookSeats(showID string, seatIDs []string) error
}

func (s *AvailabilityService) GetAvailableSeats(showID string) ([]catalog.ShowSeat, error) {
	return s.showSeatRepo.GetAvailableSeats(showID)
}

func (s *AvailabilityService) LockSeats(showID string, seatIDs []string) error {
	var successfullyLocked []string

	for _, seatID := range seatIDs {
		err := s.lockService.TryLock(seatID)
		if err != nil {
			s.releaseLocks(successfullyLocked)
			return err
		}
		
		// FIX: We MUST verify the seat is actually AVAILABLE in the database.
		// If it was already BOOKED yesterday, we should reject it!
		seat, err := s.showSeatRepo.GetByID(seatID)
		if err != nil {
			s.releaseLocks(successfullyLocked)
			_ = s.lockService.Unlock(seatID)
			return err
		}
		if seat.Status != catalog.ShowSeatStatusAvailable {
			s.releaseLocks(successfullyLocked)
			_ = s.lockService.Unlock(seatID)
			return fmt.Errorf("seat %s is not available (current status: %s)", seatID, seat.Status)
		}

		successfullyLocked = append(successfullyLocked, seatID)
	}

	err := s.showSeatRepo.UpdateStatuses(seatIDs, catalog.ShowSeatStatusLocked)
	if err != nil {
		s.releaseLocks(seatIDs)
		return err
	}

	return nil
}

func (s *AvailabilityService) ReleaseSeats(showID string, seatIDs []string) error {
	err := s.showSeatRepo.UpdateStatuses(seatIDs, catalog.ShowSeatStatusAvailable)
	s.releaseLocks(seatIDs)
	return err
}

func (s *AvailabilityService) BookSeats(showID string, seatIDs []string) error {
	err := s.showSeatRepo.UpdateStatuses(seatIDs, catalog.ShowSeatStatusBooked)
	s.releaseLocks(seatIDs)
	return err
}

func (s *AvailabilityService) releaseLocks(seatIDs []string) {
	for _, seatID := range seatIDs {
		_ = s.lockService.Unlock(seatID)
	}
}

