package postgres

import (
	"context"
	"time"
	"ticketer/internal/catalog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShowRepository struct {
	db *pgxpool.Pool
}

func NewShowRepository(db *pgxpool.Pool) catalog.ShowRepository {
	return &ShowRepository{
		db: db,
	}
}

func (r *ShowRepository) GetByID(id string) (*catalog.Show, error) {
	query := `SELECT id, movie_id, screen_id, start_time, end_time FROM shows WHERE id = $1`
	var show catalog.Show
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&show.ID,
		&show.MovieID,
		&show.ScreenID,
		&show.StartTime,
		&show.EndTime,
	)
	if err != nil {
		return nil, err
	}
	return &show, nil
}

func (r *ShowRepository) GetByMovie(movieID string) ([]catalog.Show, error) {
	query := `SELECT id, movie_id, screen_id, start_time, end_time FROM shows WHERE movie_id = $1`
	rows, err := r.db.Query(context.Background(), query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []catalog.Show
	for rows.Next() {
		var show catalog.Show
		if err := rows.Scan(
			&show.ID,
			&show.MovieID,
			&show.ScreenID,
			&show.StartTime,
			&show.EndTime,
		); err != nil {
			return nil, err
		}
		shows = append(shows, show)
	}
	return shows, nil
}

func (r *ShowRepository) GetByScreen(screenID string) ([]catalog.Show, error) {
	query := `SELECT id, movie_id, screen_id, start_time, end_time FROM shows WHERE screen_id = $1`
	rows, err := r.db.Query(context.Background(), query, screenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []catalog.Show
	for rows.Next() {
		var show catalog.Show
		if err := rows.Scan(
			&show.ID,
			&show.MovieID,
			&show.ScreenID,
			&show.StartTime,
			&show.EndTime,
		); err != nil {
			return nil, err
		}
		shows = append(shows, show)
	}
	return shows, nil
}

func (r *ShowRepository) GetByScreenAndTime(screenID string, startTime time.Time) (*catalog.Show, error) {
	query := `SELECT id, movie_id, screen_id, start_time, end_time FROM shows WHERE screen_id = $1 AND start_time = $2`
	var show catalog.Show
	err := r.db.QueryRow(context.Background(), query, screenID, startTime).Scan(
		&show.ID,
		&show.MovieID,
		&show.ScreenID,
		&show.StartTime,
		&show.EndTime,
	)
	if err != nil {
		return nil, err
	}
	return &show, nil
}

func (r *ShowRepository) GetByTheater(theaterID string) ([]catalog.Show, error) {
	query := `
		SELECT sh.id, sh.movie_id, sh.screen_id, sh.start_time, sh.end_time 
		FROM shows sh 
		JOIN screens sc ON sh.screen_id = sc.id 
		WHERE sc.theater_id = $1`
	rows, err := r.db.Query(context.Background(), query, theaterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shows []catalog.Show
	for rows.Next() {
		var show catalog.Show
		if err := rows.Scan(
			&show.ID,
			&show.MovieID,
			&show.ScreenID,
			&show.StartTime,
			&show.EndTime,
		); err != nil {
			return nil, err
		}
		shows = append(shows, show)
	}
	return shows, nil
}

func (r *ShowRepository) Save(show *catalog.Show) error {
	query := `INSERT INTO shows (id, movie_id, screen_id, start_time, end_time) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(context.Background(), query,
		show.ID,
		show.MovieID,
		show.ScreenID,
		show.StartTime,
		show.EndTime,
	)
	return err
}

func (r *ShowRepository) Update(show *catalog.Show) error {
	query := `UPDATE shows SET movie_id = $1, screen_id = $2, start_time = $3, end_time = $4 WHERE id = $5`
	_, err := r.db.Exec(context.Background(), query,
		show.MovieID,
		show.ScreenID,
		show.StartTime,
		show.EndTime,
		show.ID,
	)
	return err
}

func (r *ShowRepository) Delete(id string) error {
	query := `DELETE FROM shows WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
