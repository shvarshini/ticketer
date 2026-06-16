package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/google/login", h.handleLogin)
	mux.HandleFunc("GET /auth/google/callback", h.handleCallback)
	mux.Handle("GET /auth/me", AuthMiddleware(http.HandlerFunc(h.handleMe)))
	mux.Handle("POST /auth/logout", AuthMiddleware(http.HandlerFunc(h.handleLogout)))
}

func generateStateToken(role string) string {
	b := make([]byte, 16)
	rand.Read(b)
	csrf := base64.URLEncoding.EncodeToString(b)
	return fmt.Sprintf("%s:%s", csrf, role)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	if role == "" {
		role = "user" 
	}

	state := generateStateToken(role)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   int(time.Hour.Seconds()),
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	url := h.svc.GetLoginURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != state {
		http.Error(w, "invalid or missing state token", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	parts := strings.Split(state, ":")
	requestedRole := "user"
	if len(parts) == 2 {
		requestedRole = parts[1]
	}

	jwtToken, err := h.svc.HandleCallback(r.Context(), code, requestedRole)
	if err != nil {
		http.Error(w, "failed to authenticate: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		MaxAge:   int(time.Hour.Seconds() * 24),
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	
	http.Redirect(w, r, frontendURL, http.StatusFound)
}

func (h *Handler) handleMe(w http.ResponseWriter, r *http.Request) {
	// Assume AuthMiddleware has injected User into context
	user, ok := r.Context().Value("user").(*User)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
}
