package db

import (
	"context"
	"database/sql"
	"go-back/internal/config"
	"time"

	_ "github.com/lib/pq"
)

func ConnectToDB(cfg config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Addr)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err 
	}

	return db, nil
}