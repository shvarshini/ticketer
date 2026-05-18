package app

import (
	"fmt"
	"time"

	"ticketer/internal/catalog"

	"go.uber.org/zap"
)

func seedData(
	theaterRepo catalog.TheaterRepository,
	movieRepo catalog.MovieRepository,
	showRepo catalog.ShowRepository,
	showSeatRepo catalog.ShowSeatRepository,
	logger *zap.Logger,
) {
	logger.Info("Seeding hardcoded data for testing...")

	// 1. Create a Theater
	theater := &catalog.Theater{
		ID:       "t1",
		Name:     "PVR Cinemas",
		Location: "Downtown",
		Screens:  []catalog.Screen{},
	}
	_ = theaterRepo.Save(theater)

	// 2. Create 4 Screens, each with 4 Seats
	for i := 1; i <= 4; i++ {
		screenID := fmt.Sprintf("screen%d", i)
		screen := &catalog.Screen{
			ID:    screenID,
			Name:  fmt.Sprintf("Screen %d", i),
			Seats: []catalog.Seat{},
		}

		_ = theaterRepo.AddScreenToTheater("t1", screen)

		for j := 1; j <= 4; j++ {
			seatType := catalog.SeatTypeNormal
			if j == 4 {
				seatType = catalog.SeatTypePremium // Make the last seat premium
			}
			seat := &catalog.Seat{
				ID:       fmt.Sprintf("%s-seat%d", screenID, j),
				ScreenID: screenID,
				Row:      "A",
				Number:   j,
				Type:     seatType,
			}
			_, _ = theaterRepo.AddSeatToScreen(screenID, seat)
		}
	}

	// 3. Create a Movie to show
	movie := &catalog.Movie{
		ID:          "m1",
		Title:       "Inception",
		Description: "A thief who steals corporate secrets through the use of dream-sharing technology.",
		Duration:    148,
		ReleaseDate: time.Now(),
		Genre:       "Sci-Fi",
		BasePrice:   150.0,
	}
	_ = movieRepo.Save(movie)

	// 4. Create a Show on Screen 1
	show := &catalog.Show{
		ID:        "show1",
		MovieID:   "m1",
		ScreenID:  "screen1",
		StartTime: time.Now().Add(24 * time.Hour), // Tomorrow
		EndTime:   time.Now().Add(27 * time.Hour),
	}
	_ = showRepo.Save(show)

	// 5. Create ShowSeats (Tickets) for the 4 seats in Screen 1
	for j := 1; j <= 4; j++ {
		seatID := fmt.Sprintf("screen1-seat%d", j)
		showSeat := &catalog.ShowSeat{
			ID:     fmt.Sprintf("showseat-%d", j),
			ShowID: "show1",
			SeatID: seatID,
			Price:  150.0, // Base price, pricing service handles premium markup
			Status: catalog.ShowSeatStatusAvailable,
		}
		_ = showSeatRepo.Save(showSeat)
	}

	logger.Info("Seeding complete! You can now test bookings.")
	logger.Info("Use ShowID: 'show1'")
	logger.Info("Available SeatIDs: ['screen1-seat1', 'screen1-seat2', 'screen1-seat3', 'screen1-seat4']")
}
