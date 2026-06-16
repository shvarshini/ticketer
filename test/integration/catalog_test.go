package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ticketer/internal/catalog"
	"ticketer/test/integration/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatalogFlow(t *testing.T) {
	db := testutils.SetupTestDB(t)
	app := testutils.SetupTestApp(t, db)
	ctx := context.Background()

	// Seed Admin user
	adminID := "22222222-2222-2222-2222-222222222222"
	_, err := db.Exec(ctx, `
		INSERT INTO users (id, email, name, role, oauth_provider)
		VALUES ($1, 'admin@example.com', 'Admin User', 'admin', 'google')
	`, adminID)
	require.NoError(t, err)

	adminCookie := testutils.GenerateTestToken(adminID, "admin")

	// 1. Create a Theater
	theaterPayload := map[string]string{
		"name":     "Test Cinema",
		"location": "Downtown",
		"admin_id": adminID,
	}
	body, _ := json.Marshal(theaterPayload)
	req, _ := http.NewRequest(http.MethodPost, "/api/admin/theaters", bytes.NewBuffer(body))
	req.AddCookie(adminCookie)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	app.Mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)

	var theater catalog.Theater
	err = json.Unmarshal(rr.Body.Bytes(), &theater)
	require.NoError(t, err)
	assert.Equal(t, "Test Cinema", theater.Name)
	assert.NotEmpty(t, theater.ID)

	// 2. Add a Screen
	screenPayload := map[string]string{
		"name": "Screen 1",
	}
	body, _ = json.Marshal(screenPayload)
	req, _ = http.NewRequest(http.MethodPost, "/api/admin/theaters/"+theater.ID+"/screens", bytes.NewBuffer(body))
	req.AddCookie(adminCookie)
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	app.Mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)

	var screen catalog.Screen
	err = json.Unmarshal(rr.Body.Bytes(), &screen)
	require.NoError(t, err)
	assert.Equal(t, "Screen 1", screen.Name)
	assert.NotEmpty(t, screen.ID)

	// 3. Add a Movie
	moviePayload := map[string]interface{}{
		"title":       "Inception",
		"description": "Dream within a dream",
		"duration":    148,
		"language":    "English",
		"genre":       "Sci-Fi",
	}
	body, _ = json.Marshal(moviePayload)
	req, _ = http.NewRequest(http.MethodPost, "/api/admin/movies", bytes.NewBuffer(body))
	req.AddCookie(adminCookie)
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	app.Mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)

	var movie catalog.Movie
	err = json.Unmarshal(rr.Body.Bytes(), &movie)
	require.NoError(t, err)
	assert.Equal(t, "Inception", movie.Title)
	assert.NotEmpty(t, movie.ID)
}
