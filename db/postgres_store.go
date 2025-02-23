package db

import (
	"database/sql"
	"fmt"
	"github.com/finlleyl/shorty/internal/app"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (p *PostgresStore) Save(shortURL, originalURL string) int {
	var id int
	err := p.db.QueryRow(
		"INSERT INTO urls (short_url, original_url) VALUES ($1, $2) RETURNING id", shortURL, originalURL).Scan(&id)
	if err != nil {
		fmt.Println("Error inserting URL:", err)
		return 0
	}

	return id
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
