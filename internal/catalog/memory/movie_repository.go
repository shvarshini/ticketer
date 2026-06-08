package memory

import (
	"errors"
	"sync"
	"ticketer/internal/catalog"
)

type MovieRepository struct {
	mu     sync.RWMutex
	movies map[string]*catalog.Movie
}


func NewMovieRepository() *MovieRepository {
	return &MovieRepository{
		movies: make(map[string]*catalog.Movie),
	}
}

func (r *MovieRepository) GetByID(id string) (*catalog.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	movie, ok := r.movies[id]
	if !ok {
		return nil, errors.New("movie not found")
	}
	return movie, nil
}

func (r *MovieRepository) List() ([]catalog.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	movies := make([]catalog.Movie, 0, len(r.movies))
	for _, m := range r.movies {
		movies = append(movies, *m)
	}
	return movies, nil
}

func (r *MovieRepository) Save(movie *catalog.Movie) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if movie.ID == "" {
		return errors.New("movie ID is required")
	}
	r.movies[movie.ID] = movie
	return nil
}

func (r *MovieRepository) Update(movie *catalog.Movie) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.movies[movie.ID]; !ok {
		return errors.New("movie not found")
	}
	r.movies[movie.ID] = movie
	return nil
}

func (r *MovieRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.movies[id]; !ok {
		return nil
	}
	delete(r.movies, id)
	return nil
}
