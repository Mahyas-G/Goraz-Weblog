package db

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
)

func Connect(databaseURL string) (*sqlx.DB, error) {
	conn, err := sqlx.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("db: open failed: %w", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("db: ping failed: %w", err)
	}

	return conn, nil
}

func Ping(ctx context.Context, conn *sqlx.DB) error {
	var result int
	return conn.QueryRowContext(ctx, "SELECT 1").Scan(&result)
}
