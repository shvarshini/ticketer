package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ticketer/internal/booking"
	"ticketer/test/integration/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingFlow(t *testing.T) {
	db := testutils.SetupTestDB(t)
	app := testutils.SetupTestApp(t, db)
	ctx := context.Background()

	// 1. Seed the Database
	theaterID := uuid.New()
	adminID := uuid.New()
	screenID := uuid.New()
	movieID := uuid.New()
	showID := uuid.New()
	seatID1 := uuid.New()
	seatID2 := uuid.New()
	userID := uuid.New()

	// Insert Admin and User
	_, err := db.Exec(ctx, `INSERT INTO users (id, email, name, role, oauth_provider) VALUES ($1, 'admin@test.com', 'Admin', 'admin', 'google')`, adminID)
	require.NoError(t, err)
	_, err = db.Exec(ctx, `INSERT INTO users (id, email, name, role, oauth_provider) VALUES ($1, 'user@test.com', 'User', 'user', 'google')`, userID)
	require.NoError(t, err)

	// Insert Theater & Screen
	_, err = db.Exec(ctx, `INSERT INTO theaters (id, name, location, admin_id) VALUES ($1, 'Booking Theater', 'City', $2)`, theaterID, adminID)
	require.NoError(t, err)
	_, err = db.Exec(ctx, `INSERT INTO screens (id, theater_id, name) VALUES ($1, $2, 'Screen 1')`, screenID, theaterID)
	require.NoError(t, err)

	_, err = db.Exec(ctx, `INSERT INTO movies (id, title, description, duration, release_date, genre, base_price) VALUES ($1, 'Test Movie', 'Desc', 120, '2025-01-01', 'Action', 15.00)`, movieID)
	require.NoError(t, err)
	_, err = db.Exec(ctx, `INSERT INTO shows (id, movie_id, screen_id, start_time, end_time) VALUES ($1, $2, $3, $4, $5)`, showID, movieID, screenID, time.Now().Add(time.Hour), time.Now().Add(3*time.Hour))
	require.NoError(t, err)

	// Insert Seats
	_, err = db.Exec(ctx, `INSERT INTO seats (id, screen_id, row, number, type) VALUES ($1, $2, 'A', 1, 'standard')`, seatID1, screenID)
	require.NoError(t, err)
	_, err = db.Exec(ctx, `INSERT INTO seats (id, screen_id, row, number, type) VALUES ($1, $2, 'A', 2, 'premium')`, seatID2, screenID)
	require.NoError(t, err)

	// Insert ShowSeats
	showSeatID1 := uuid.New()
	showSeatID2 := uuid.New()
	_, err = db.Exec(ctx, `INSERT INTO show_seats (id, show_id, seat_id, status) VALUES ($1, $2, $3, 'AVAILABLE')`, showSeatID1, showID, seatID1)
	require.NoError(t, err)
	_, err = db.Exec(ctx, `INSERT INTO show_seats (id, show_id, seat_id, status) VALUES ($1, $2, $3, 'AVAILABLE')`, showSeatID2, showID, seatID2)
	require.NoError(t, err)

	userCookie := testutils.GenerateTestToken(userID.String(), "user")

	// 2. Initiate Booking
	bookingPayload := map[string]interface{}{
		"user_id": userID.String(),
		"show_id": showID.String(),
		"seat_ids": []string{
			showSeatID1.String(),
			showSeatID2.String(),
		},
	}
	body, _ := json.Marshal(bookingPayload)
	req, _ := http.NewRequest(http.MethodPost, "/api/bookings", bytes.NewBuffer(body))
	req.AddCookie(userCookie)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	app.Mux.ServeHTTP(rr, req)
	
	// Print response body for debugging if it fails
	require.Equal(t, http.StatusCreated, rr.Code, "Booking failed: %s", rr.Body.String())

	var b booking.Booking
	err = json.Unmarshal(rr.Body.Bytes(), &b)
	require.NoError(t, err)
	assert.Equal(t, booking.BookingStatusPending, b.Status)
	assert.Len(t, b.SeatIDs, 2)
}
