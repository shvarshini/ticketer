package availability

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	availabilityService Service
}

func NewHandler(as Service) *Handler {
	return &Handler{
		availabilityService: as,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/shows/{id}/seats", h.getAvailableSeats)
}

func (h *Handler) getAvailableSeats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	seats, err := h.availabilityService.GetAvailableSeats(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, seats)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
