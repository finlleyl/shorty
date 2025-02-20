package db

import (
	"context"
	"fmt"
	"github.com/finlleyl/shorty/internal/config"
	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func InitDB(con *config.Config) {
	databaseURL := con.D.Address

	var err error
	DB, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		panic(err)
	}
}

func PingDB(ctx context.Context) error {

	err := DB.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database is not reachable: %w", err)
	}
	return nil
}

func CloseDB() {
	DB.Close(context.Background())
}
