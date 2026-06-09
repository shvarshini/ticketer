package postgres

import (
	"context"
	"ticketer/internal/catalog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) catalog.MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) GetByID(id string) (*catalog.Movie, error) {
	query := `SELECT id, title, description, duration, release_date, genre, base_price FROM movies WHERE id = $1`
	var movie catalog.Movie
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Duration,
		&movie.ReleaseDate,
		&movie.Genre,
		&movie.BasePrice,
	)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) List() ([]catalog.Movie, error) {
	query := `SELECT id, title, description, duration, release_date, genre, base_price FROM movies`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []catalog.Movie
	for rows.Next() {
		var movie catalog.Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Duration,
			&movie.ReleaseDate,
			&movie.Genre,
			&movie.BasePrice,
		); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (r *MovieRepository) Save(movie *catalog.Movie) error {
	query := `INSERT INTO movies (id, title, description, duration, release_date, genre, base_price) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(context.Background(), query,
		movie.ID,
		movie.Title,
		movie.Description,
		movie.Duration,
		movie.ReleaseDate,
		movie.Genre,
		movie.BasePrice,
	)
	return err
}

func (r *MovieRepository) Update(movie *catalog.Movie) error {
	query := `UPDATE movies SET title = $1, description = $2, duration = $3, release_date = $4, genre = $5, base_price = $6 WHERE id = $7`
	_, err := r.db.Exec(context.Background(), query,
		movie.Title,
		movie.Description,
		movie.Duration,
		movie.ReleaseDate,
		movie.Genre,
		movie.BasePrice,
		movie.ID,
	)
	return err
}

func (r *MovieRepository) Delete(id string) error {
	query := `DELETE FROM movies WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
