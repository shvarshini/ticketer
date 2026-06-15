package postgres

import (
	"context"
	"errors"

	"ticketer/internal/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) auth.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*auth.User, error) {
	query := `SELECT id, email, name, roles, oauth_provider, created_at, updated_at FROM users WHERE id = $1`
	var u auth.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Name, &u.Roles, &u.OAuthProvider, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	query := `SELECT id, email, name, roles, oauth_provider, created_at, updated_at FROM users WHERE email = $1`
	var u auth.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.Name, &u.Roles, &u.OAuthProvider, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(ctx context.Context, user *auth.User) error {
	query := `
		INSERT INTO users (id, email, name, roles, oauth_provider, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.Name, user.Roles, user.OAuthProvider, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *auth.User) error {
	query := `
		UPDATE users 
		SET name = $1, roles = $2, updated_at = CURRENT_TIMESTAMP
		WHERE email = $3
	`
	_, err := r.db.Exec(ctx, query, user.Name, user.Roles, user.Email)
	return err
}
