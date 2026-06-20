package catalog

import (
	"errors"

	"github.com/google/uuid"
)

type TheaterService struct {
	theaterRepo TheaterRepository
}

func NewTheaterService(tr TheaterRepository) *TheaterService {
	return &TheaterService{theaterRepo: tr}
}

func (s *TheaterService) AddTheater(theater *Theater) error {
	theater.ID = uuid.New().String()
	return s.theaterRepo.Save(theater)
}

func (s *TheaterService) GetTheater(id string) (*Theater, error) {
	return s.theaterRepo.GetByID(id)
}

func (s *TheaterService) ListTheatersByAdmin(adminID string) ([]Theater, error) {
	return s.theaterRepo.GetByAdminID(adminID)
}

func (s *TheaterService) UpdateTheater(theater *Theater) error {
	_, err := s.theaterRepo.GetByID(theater.ID)
	if err != nil {
		return err
	}
	return s.theaterRepo.Update(theater)
}

func (s *TheaterService) DeleteTheater(id string) error {
	return s.theaterRepo.Delete(id)
}

func (s *TheaterService) AddScreenToTheater(theaterID string, screen *Screen) error {
	screen.ID = uuid.New().String()
	screen.TheaterID = theaterID
	return s.theaterRepo.AddScreenToTheater(theaterID, screen)
}

func (s *TheaterService) UpdateScreen(screen *Screen) error {
	return s.theaterRepo.UpdateScreen(screen)
}

func (s *TheaterService) DeleteScreen(theaterID, screenID string) error {
	return s.theaterRepo.DeleteScreen(theaterID, screenID)
}

func (s *TheaterService) GetScreens(theaterID string) ([]Screen, error) {
	return s.theaterRepo.GetScreens(theaterID)
}

func (s *TheaterService) GetScreen(screenID string) (*Screen, error) {
	return s.theaterRepo.GetScreen(screenID)
}

func (s *TheaterService) AddSeatToScreen(screenID string, seat *Seat) error {
	seat.ID = uuid.New().String()
	seat.ScreenID = screenID
	return s.theaterRepo.AddSeatToScreen(screenID, seat)
}

func (s *TheaterService) UpdateSeat(seat *Seat) error {
	return s.theaterRepo.UpdateSeat(seat)
}

func (s *TheaterService) DeleteSeat(screenID, seatID string) error {
	return s.theaterRepo.DeleteSeat(screenID, seatID)
}

func (s *TheaterService) GetSeats(screenID string) ([]Seat, error) {
	return s.theaterRepo.GetSeats(screenID)
}

type MovieService struct {
	movieRepo MovieRepository
}

func NewMovieService(mr MovieRepository) *MovieService {
	return &MovieService{movieRepo: mr}
}

func (s *MovieService) AddMovie(movie *Movie) error {
	movie.ID = uuid.New().String()
	return s.movieRepo.Save(movie)
}

func (s *MovieService) GetMovie(id string) (*Movie, error) {
	return s.movieRepo.GetByID(id)
}

func (s *MovieService) ListMovies() ([]Movie, error) {
	return s.movieRepo.List()
}

func (s *MovieService) UpdateMovie(movie *Movie) error {
	if _, err := s.movieRepo.GetByID(movie.ID); err != nil {
		return errors.New("movie doesn't exist")
	}
	return s.movieRepo.Update(movie)
}

func (s *MovieService) DeleteMovie(id string) error {
	return s.movieRepo.Delete(id)
}

type ShowService struct {
	showRepo     ShowRepository
	showSeatRepo ShowSeatRepository
	theaterRepo  TheaterRepository
}

func NewShowService(sr ShowRepository, ssr ShowSeatRepository, tr TheaterRepository) *ShowService {
	return &ShowService{
		showRepo:     sr,
		showSeatRepo: ssr,
		theaterRepo:  tr,
	}
}

func (s *ShowService) AddShow(show *Show) error {
	existingShow, err := s.showRepo.GetByScreenAndTime(show.ScreenID, show.StartTime)
	if err == nil {
		*show = *existingShow
		return nil	
	}

	show.ID = uuid.New().String()
	seats, err :=s.theaterRepo.GetSeats(show.ScreenID)
	if err != nil {
		return err
	}
	err = s.showRepo.Save(show)
	if err != nil {
		return err
	}

	for _, seat := range seats {
		showSeat := &ShowSeat{
			ID:     uuid.New().String(),
			ShowID: show.ID,
			SeatID: seat.ID,
			Status: ShowSeatStatusAvailable,
		}
		if err := s.showSeatRepo.Save(showSeat); err != nil {
			return err
		}
	}
	return nil
}

func (s *ShowService) GetShow(id string) (*Show, error) {
	return s.showRepo.GetByID(id)
}

func (s *ShowService) GetShowsByMovie(movieID string) ([]Show, error) {
	return s.showRepo.GetByMovie(movieID)
}

func (s *ShowService) GetShowsByTheater(theaterID string) ([]Show, error) {
	return s.showRepo.GetByTheater(theaterID)
}

func (s *ShowService) UpdateShow(show *Show) error {
	if _, err := s.showRepo.GetByID(show.ID); err != nil {
		return errors.New("show doesn't exist")
	}
	return s.showRepo.Update(show)
}

func (s *ShowService) DeleteShow(id string) error {
	return s.showRepo.Delete(id)
}
