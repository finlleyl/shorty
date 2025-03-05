package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/finlleyl/shorty/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDB(con *config.Config) error {
	databaseURL := con.D.Address

	var err error
	DB, err = sql.Open("pgx", databaseURL)
	if err != nil {
		panic(err)
	}

	query := `
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
	short_url TEXT UNIQUE NOT NULL,
	original_url TEXT UNIQUE NOT NULL,
    user_id TEXT NOT NULL,
    deleted_flag BOOL NOT NULL DEFAULT FALSE
);`
	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	return nil
}

func PingDB() error {

	err := DB.PingContext(context.Background())
	if err != nil {
		return fmt.Errorf("database is not reachable: %w", err)
	}
	return nil
}

func CloseDB() {
	DB.Close()
}
