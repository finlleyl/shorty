package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/finlleyl/shorty/internal/app"
	"github.com/finlleyl/shorty/internal/apperrors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (p *PostgresStore) Save(shortURL, originalURL string) (int, error) {
	var id int
	err := p.db.QueryRow(
		"INSERT INTO urls (short_url, original_url) VALUES ($1, $2) RETURNING id", shortURL, originalURL).Scan(&id)

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
	rows, err := p.db.Query("SELECT id, short_url, original_url FROM urls")
	if err != nil {
		fmt.Println("Error fetching URLs:", err)
		return nil
	}
	defer rows.Close()

	var results []app.ShortResult
	for rows.Next() {
		var r app.ShortResult
		if err := rows.Scan(&r.ID, &r.ShortURL, &r.OriginalURL); err != nil {
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
