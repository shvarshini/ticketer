package catalog

import (
	"encoding/json"
	"net/http"

	"ticketer/internal/auth"
)

type Handler struct {
	theaterService *TheaterService
	movieService   *MovieService
	showService    *ShowService
}

func NewHandler(ts *TheaterService, ms *MovieService, ss *ShowService) *Handler {
	return &Handler{
		theaterService: ts,
		movieService:   ms,
		showService:    ss,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Helper for admin routes
	adminOnly := func(next http.HandlerFunc) http.Handler {
		return auth.AuthMiddleware(auth.RoleMiddleware("admin")(next))
	}

	// Theaters
	mux.Handle("POST /api/admin/theaters", adminOnly(h.addTheater))
	mux.Handle("GET /api/admin/theaters", adminOnly(h.listTheatersByAdmin))
	mux.HandleFunc("GET /api/theaters/{id}", h.getTheater)
	mux.Handle("PUT /api/admin/theaters/{id}", adminOnly(h.updateTheater))
	mux.Handle("DELETE /api/admin/theaters/{id}", adminOnly(h.deleteTheater))

	// Screens
	mux.Handle("POST /api/admin/theaters/{id}/screens", adminOnly(h.addScreen))
	mux.HandleFunc("GET /api/theaters/{id}/screens", h.getScreens)
	mux.Handle("PUT /api/admin/screens/{screen_id}", adminOnly(h.updateScreen))
	mux.Handle("DELETE /api/admin/theaters/{id}/screens/{screen_id}", adminOnly(h.deleteScreen))
	mux.HandleFunc("GET /api/screens/{id}", h.getScreen)

	// Seat routes
	mux.Handle("POST /api/admin/screens/{screen_id}/seats", adminOnly(h.addSeat))
	mux.HandleFunc("GET /api/screens/{screen_id}/seats", h.getSeats)
	mux.Handle("PUT /api/admin/seats/{seat_id}", adminOnly(h.updateSeat))
	mux.Handle("DELETE /api/admin/screens/{screen_id}/seats/{seat_id}", adminOnly(h.deleteSeat))

	// Movies
	mux.Handle("POST /api/admin/movies", adminOnly(h.addMovie))
	mux.HandleFunc("GET /api/movies", h.listMovies)
	mux.HandleFunc("GET /api/movies/{id}", h.getMovie)
	mux.Handle("PUT /api/admin/movies/{id}", adminOnly(h.updateMovie))
	mux.Handle("DELETE /api/admin/movies/{id}", adminOnly(h.deleteMovie))

	// Shows
	mux.Handle("POST /api/admin/shows", adminOnly(h.addShow))
	mux.HandleFunc("GET /api/shows/{id}", h.getShow)
	mux.HandleFunc("GET /api/movies/{id}/shows", h.getShowsByMovie)
	mux.HandleFunc("GET /api/theaters/{id}/shows", h.getShowsByTheater)
	mux.Handle("PUT /api/admin/shows/{id}", adminOnly(h.updateShow))
	mux.Handle("DELETE /api/admin/shows/{id}", adminOnly(h.deleteShow))
}

// helper for JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// helper for error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// --- Theater Handlers ---

func (h *Handler) addTheater(w http.ResponseWriter, r *http.Request) {
	var req Theater
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.theaterService.AddTheater(&req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, req)
}

func (h *Handler) listTheatersByAdmin(w http.ResponseWriter, r *http.Request) {
	adminID := r.URL.Query().Get("admin_id")
	if adminID == "" {
		respondError(w, http.StatusBadRequest, "admin_id query parameter is required")
		return
	}

	theaters, err := h.theaterService.ListTheatersByAdmin(adminID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, theaters)
}

func (h *Handler) getTheater(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	th, err := h.theaterService.GetTheater(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, th)
}

func (h *Handler) updateTheater(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req Theater
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id
	if err := h.theaterService.UpdateTheater(&req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, req)
}

func (h *Handler) deleteTheater(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.theaterService.DeleteTheater(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Screen Handlers ---

func (h *Handler) addScreen(w http.ResponseWriter, r *http.Request) {
	theaterID := r.PathValue("id")
	var req Screen
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.theaterService.AddScreenToTheater(theaterID, &req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, req)
}

func (h *Handler) getScreens(w http.ResponseWriter, r *http.Request) {
	theaterID := r.PathValue("id")
	screens, err := h.theaterService.GetScreens(theaterID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, screens)
}

func (h *Handler) getScreen(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	// We need GetScreen in TheaterService. It's not there yet, let me check.
	// Actually, wait, `TheaterService` has no `GetScreen`. Oh wait, `s.theaterRepo.GetScreen(screenID)` exists! 
	// I need to add `GetScreen(id string) (*Screen, error)` to `TheaterService` in `service.go`. Wait! 
	// I didn't add it to `service.go`. Let me just add it to service.go, but for now I will write the handler.
	screen, err := h.theaterService.GetScreen(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, screen)
}
func (h *Handler) updateScreen(w http.ResponseWriter, r *http.Request) {
	screenID := r.PathValue("screen_id")
	var req Screen
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = screenID
	if err := h.theaterService.UpdateScreen(&req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, req)
}

func (h *Handler) deleteScreen(w http.ResponseWriter, r *http.Request) {
	theaterID := r.PathValue("id")
	screenID := r.PathValue("screen_id")
	if err := h.theaterService.DeleteScreen(theaterID, screenID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Seat Handlers ---

func (h *Handler) addSeat(w http.ResponseWriter, r *http.Request) {
	screenID := r.PathValue("screen_id")
	var req Seat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if _, err := h.theaterService.AddSeatToScreen(screenID, &req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, req)
}

func (h *Handler) getSeats(w http.ResponseWriter, r *http.Request) {
	screenID := r.PathValue("screen_id")
	seats, err := h.theaterService.GetSeats(screenID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, seats)
}

func (h *Handler) updateSeat(w http.ResponseWriter, r *http.Request) {
	seatID := r.PathValue("seat_id")
	var req Seat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = seatID
	if err := h.theaterService.UpdateSeat(&req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, req)
}

func (h *Handler) deleteSeat(w http.ResponseWriter, r *http.Request) {
	screenID := r.PathValue("screen_id")
	seatID := r.PathValue("seat_id")
	if err := h.theaterService.DeleteSeat(screenID, seatID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Movie Handlers ---

func (h *Handler) addMovie(w http.ResponseWriter, r *http.Request) {
	var req Movie
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.movieService.AddMovie(&req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, req)
}

func (h *Handler) listMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.movieService.ListMovies()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, movies)
}

func (h *Handler) getMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	movie, err := h.movieService.GetMovie(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, movie)
}

func (h *Handler) updateMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req Movie
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id
	if err := h.movieService.UpdateMovie(&req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, req)
}

func (h *Handler) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.movieService.DeleteMovie(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Show Handlers ---

func (h *Handler) addShow(w http.ResponseWriter, r *http.Request) {
	var req Show
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.showService.AddShow(&req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, req)
}

func (h *Handler) getShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	show, err := h.showService.GetShow(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, show)
}

func (h *Handler) getShowsByMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	shows, err := h.showService.GetShowsByMovie(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, shows)
}

func (h *Handler) getShowsByTheater(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	shows, err := h.showService.GetShowsByTheater(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, shows)
}

func (h *Handler) updateShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req Show
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id
	if err := h.showService.UpdateShow(&req); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, req)
}

func (h *Handler) deleteShow(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.showService.DeleteShow(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


