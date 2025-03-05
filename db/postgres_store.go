package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/apperrors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"sync"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (p *PostgresStore) Save(shortURL, originalURL, userID string) (int, error) {
	var id int
	err := p.db.QueryRow(
		"INSERT INTO urls (short_url, original_url, user_id) VALUES ($1, $2, $3) RETURNING id", shortURL, originalURL, userID).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			var existingShortURL string
			queryErr := p.db.QueryRow(
				"SELECT short_url FROM urls WHERE original_url = $1",
				originalURL,
			).Scan(&existingShortURL)

			if queryErr != nil {
				return 0, queryErr
			}

			return 0, apperrors.NewConflictError(existingShortURL)
		}

		return 0, err
	}

	return id, nil
}

func (p *PostgresStore) Get(id string) (string, bool) {
	var originalURL string
	err := p.db.QueryRow(""+
		"SELECT original_url FROM urls where short_url = $1", id).Scan(&originalURL)

	if err != nil {
		return "", false
	}

	return originalURL, true
}

func (p *PostgresStore) GetAll() []app.ShortResult {
	rows, err := p.db.Query("SELECT id, short_url, original_url, user_id FROM urls")
	if err != nil {
		fmt.Println("Error fetching URLs:", err)
		return nil
	}
	defer rows.Close()

	var results []app.ShortResult
	for rows.Next() {
		var r app.ShortResult
		if err := rows.Scan(&r.ID, &r.ShortURL, &r.OriginalURL, &r.UserID); err != nil {
			continue
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
		return nil
	}

	return results
}

func (p *PostgresStore) GetByUserID(userID string) ([]app.ShortResult, error) {

	rows, err := p.db.Query("SELECT id, short_url, original_url, user_id FROM urls WHERE user_id = $1", userID)
	if err != nil {
		fmt.Println("Error fetching URLs:", err)
		return nil, err
	}
	defer rows.Close()
	var results []app.ShortResult
	for rows.Next() {
		var r app.ShortResult
		if err := rows.Scan(&r.ID, &r.ShortURL, &r.OriginalURL, &r.UserID); err != nil {
			continue
		}
		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
		return nil, err
	}

	return results, nil
}

func (p *PostgresStore) BatchDelete(urls []string, userID string) error {
	const batchSize = 100
	errs := make(chan error, len(urls)/batchSize+1)
	var wg sync.WaitGroup

	for start := 0; start < len(urls); start += batchSize {
		end := start + batchSize
		if end > len(urls) {
			end = len(urls)
		}
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			query := `UPDATE urls SET deleted_flag = TRUE WHERE user_id = $1 AND short_url = ANY($2)`
			if _, err := p.db.Exec(query, userID, pq.Array(batch)); err != nil {
				errs <- err
			}
		}(urls[start:end])
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
