package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/finlleyl/shorty/internal/config"
	_ "github.com/jackc/pgx/v5"
)

var DB *sql.DB

func InitDB(con *config.Config) {
	databaseURL := con.D.Address

	var err error
	DB, err = sql.Open("pgx", databaseURL)
	if err != nil {
		panic(err)
	}
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
