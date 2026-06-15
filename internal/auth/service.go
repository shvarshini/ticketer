package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type service struct {
	repo         UserRepository
	oauthConfig  *oauth2.Config
	jwtSecret    []byte
}

func NewService(repo UserRepository) Service {
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/auth/google/callback"
	}

	conf := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "super_secret_key_for_development"
	}

	return &service{
		repo:        repo,
		oauthConfig: conf,
		jwtSecret:   []byte(secret),
	}
}

func (s *service) GetLoginURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *service) HandleCallback(ctx context.Context, code string, requestedRole string) (string, error) {
	
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}

	client := s.oauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	defer resp.Body.Close()
	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return "", fmt.Errorf("failed to decode user info: %w", err)
	}

	user, err := s.repo.FindByEmail(ctx, googleUser.Email)
	if err != nil {
		return "", fmt.Errorf("failed to check existing user: %w", err)
	}

	if user == nil {
		var roles []string
		switch requestedRole {
		case "admin":
			roles = []string{"admin"}
		case "both":
			roles = []string{"user", "admin"}
		default:
			roles = []string{"user"} 
		}

		user = &User{
			ID:            uuid.New(),
			Email:         googleUser.Email,
			Name:          googleUser.Name,
			Roles:         roles,
			OAuthProvider: "google",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := s.repo.Create(ctx, user); err != nil {
			return "", fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		user.Name = googleUser.Name
		switch requestedRole {
		case "admin":
			user.Roles = []string{"admin"}
		case "both":
			user.Roles = []string{"user", "admin"}
		default:
			user.Roles = []string{"user"}
		}
		if err := s.repo.Update(ctx, user); err != nil {
			return "", fmt.Errorf("failed to update user: %w", err)
		}
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID.String(),
		"roles": user.Roles,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), 
		"iat":   time.Now().Unix(),
	})

	tokenString, err := jwtToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *service) GetProfile(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
