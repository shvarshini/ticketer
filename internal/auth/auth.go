package auth

import (
	"context"
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	Roles         []string  `json:"roles"`
	OAuthProvider string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type Service interface {
	GetLoginURL(state string) string
	HandleCallback(ctx context.Context, code string, requestedRole string) (string, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
}
