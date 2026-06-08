package memory

import (
	"fmt"
	"sync"
	"ticketer/internal/catalog"
)

type TheaterRepository struct{
	mu     sync.RWMutex
	theater map[string]*catalog.Theater
	screenTheaterMap map[string]string
}


func NewTheaterRepository() *TheaterRepository {
	return &TheaterRepository{
		theater:          make(map[string]*catalog.Theater),
		screenTheaterMap: make(map[string]string),
	}
}

func (tRepo *TheaterRepository) GetByID(id string) (*catalog.Theater, error) {
	tRepo.mu.RLock()
	defer tRepo.mu.RUnlock()

	if val, ok := tRepo.theater[id]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("no theater found for ID: %s", id)
}

func (tRepo *TheaterRepository) GetScreen(screenID string) (*catalog.Screen, error) {
	tRepo.mu.RLock()
	defer tRepo.mu.RUnlock()

	theaterID, ok := tRepo.screenTheaterMap[screenID]
	if !ok {
		return nil, fmt.Errorf("no theater found for screen ID: %s", screenID)
	}

	theater, ok := tRepo.theater[theaterID]
	if !ok {
		return nil, fmt.Errorf("theater %s not found for screen %s", theaterID, screenID)
	}

	for _, sc := range theater.Screens {
		if sc.ID == screenID {
			return &sc, nil
		}
	}

	return nil, fmt.Errorf("no screen with ID %s found in theater %s", screenID, theaterID)
}

func (tRepo *TheaterRepository) List() ([]catalog.Theater, error) {
	tRepo.mu.RLock()
	defer tRepo.mu.RUnlock()

	theaterList := make([]catalog.Theater, 0, len(tRepo.theater))
	for _, v := range tRepo.theater {
		theaterList = append(theaterList, *v)
	}

	return theaterList, nil
}

func (tRepo *TheaterRepository) Save(theater *catalog.Theater) error {
	tRepo.mu.Lock()
	defer tRepo.mu.Unlock()

	tRepo.theater[theater.ID] = theater
	for _, screen := range theater.Screens {
		tRepo.screenTheaterMap[screen.ID] = theater.ID
	}
	return nil
}

func (tRepo *TheaterRepository) AddScreenToTheater(theaterID string, screen *catalog.Screen) error {
	tRepo.mu.Lock()
	defer tRepo.mu.Unlock()

	theater, ok := tRepo.theater[theaterID]
	if !ok {
		return fmt.Errorf("theater not found for ID: %s", theaterID)
	}
	screen.TheaterID = theaterID
	theater.Screens = append(theater.Screens, *screen)
	tRepo.screenTheaterMap[screen.ID] = theaterID
	return nil
}

func (tRepo *TheaterRepository) AddSeatToScreen(screenID string, seat *catalog.Seat) (*catalog.Screen, error) {
	tRepo.mu.Lock()
	defer tRepo.mu.Unlock()

	theaterID, ok := tRepo.screenTheaterMap[screenID]
	if !ok {
		return nil, fmt.Errorf("no theater found for screen ID: %s", screenID)
	}

	theater, ok := tRepo.theater[theaterID]
	if !ok {
		return nil, fmt.Errorf("theater %s not found for screen %s", theaterID, screenID)
	}

	for i := range theater.Screens {
		if theater.Screens[i].ID == screenID {
			theater.Screens[i].Seats = append(theater.Screens[i].Seats, *seat)
			return &theater.Screens[i], nil
		}
	}

	return nil, fmt.Errorf("screen with ID %s not found in theater %s", screenID, theaterID)
}

func (tRepo *TheaterRepository) GetByAdminID(adminID string) ([]catalog.Theater, error) {
	tRepo.mu.RLock()
	defer tRepo.mu.RUnlock()

	var result []catalog.Theater
	for _, th := range tRepo.theater {
		if th.AdminID == adminID {
			result = append(result, *th)
		}
	}
	return result, nil
}

func (tRepo *TheaterRepository) Update(theater *catalog.Theater) error {
	tRepo.mu.Lock()
	defer tRepo.mu.Unlock()

	if _, ok := tRepo.theater[theater.ID]; !ok {
		return fmt.Errorf("theater not found for ID: %s", theater.ID)
	}
	tRepo.theater[theater.ID] = theater
	return nil
}

func (tRepo *TheaterRepository) Delete(id string) error {
	tRepo.mu.Lock()
	defer tRepo.mu.Unlock()

	theater, ok := tRepo.theater[id]
	if !ok {
		return nil
	}

	for _, screen := range theater.Screens {
		delete(tRepo.screenTheaterMap, screen.ID)
	}
	delete(tRepo.theater, id)
	return nil
}

func (tRepo *TheaterRepository) UpdateScreen(screen *catalog.Screen) error {
	tRepo.mu.Lock()
	defer tRepo.mu.Unlock()

	theaterID, ok := tRepo.screenTheaterMap[screen.ID]
	if !ok {
		return fmt.Errorf("no theater found for screen ID: %s", screen.ID)
	}

	theater, ok := tRepo.theater[theaterID]
	if !ok {
		return fmt.Errorf("theater %s not found for screen %s", theaterID, screen.ID)
	}

	for i := range theater.Screens {
		if theater.Screens[i].ID == screen.ID {
			screen.Seats = theater.Screens[i].Seats
			theater.Screens[i] = *screen
			return nil
		}
	}
	return fmt.Errorf("screen with ID %s not found in theater %s", screen.ID, theaterID)
}

func (tRepo *TheaterRepository) DeleteScreen(theaterID, screenID string) error {
	tRepo.mu.Lock()
	defer tRepo.mu.Unlock()

	theater, ok := tRepo.theater[theaterID]
	if !ok {
		return nil
	}

	for i := range theater.Screens {
		if theater.Screens[i].ID == screenID {
			theater.Screens = append(theater.Screens[:i], theater.Screens[i+1:]...)
			delete(tRepo.screenTheaterMap, screenID)
			return nil
		}
	}
	return nil
}

func (tRepo *TheaterRepository) GetSeats(screenID string) ([]catalog.Seat, error) {
	tRepo.mu.RLock()
	defer tRepo.mu.RUnlock()

	theaterID, ok := tRepo.screenTheaterMap[screenID]
	if !ok {
		return nil, fmt.Errorf("no theater found for screen ID: %s", screenID)
	}

	theater, ok := tRepo.theater[theaterID]
	if !ok {
		return nil, fmt.Errorf("theater %s not found for screen %s", theaterID, screenID)
	}

	for _, sc := range theater.Screens {
		if sc.ID == screenID {
			return sc.Seats, nil
		}
	}
	return nil, fmt.Errorf("screen with ID %s not found in theater %s", screenID, theaterID)
}

func (tRepo *TheaterRepository) UpdateSeat(seat *catalog.Seat) error {
	tRepo.mu.Lock()
	defer tRepo.mu.Unlock()

	theaterID, ok := tRepo.screenTheaterMap[seat.ScreenID]
	if !ok {
		return fmt.Errorf("no theater found for screen ID: %s", seat.ScreenID)
	}

	theater, ok := tRepo.theater[theaterID]
	if !ok {
		return fmt.Errorf("theater %s not found for screen %s", theaterID, seat.ScreenID)
	}

	for i := range theater.Screens {
		if theater.Screens[i].ID == seat.ScreenID {
			for j := range theater.Screens[i].Seats {
				if theater.Screens[i].Seats[j].ID == seat.ID {
					theater.Screens[i].Seats[j] = *seat
					return nil
				}
			}
			return fmt.Errorf("seat with ID %s not found in screen %s", seat.ID, seat.ScreenID)
		}
	}
	return fmt.Errorf("screen with ID %s not found in theater %s", seat.ScreenID, theaterID)
}

func (tRepo *TheaterRepository) DeleteSeat(screenID, seatID string) error {
	tRepo.mu.Lock()
	defer tRepo.mu.Unlock()

	theaterID, ok := tRepo.screenTheaterMap[screenID]
	if !ok {
		return fmt.Errorf("no theater found for screen ID: %s", screenID)
	}

	theater, ok := tRepo.theater[theaterID]
	if !ok {
		return fmt.Errorf("theater %s not found for screen %s", theaterID, screenID)
	}

	for i := range theater.Screens {
		if theater.Screens[i].ID == screenID {
			for j := range theater.Screens[i].Seats {
				if theater.Screens[i].Seats[j].ID == seatID {
					theater.Screens[i].Seats = append(theater.Screens[i].Seats[:j], theater.Screens[i].Seats[j+1:]...)
					return nil
				}
			}
			return nil
		}
	}
	return fmt.Errorf("screen with ID %s not found in theater %s", screenID, theaterID)
}