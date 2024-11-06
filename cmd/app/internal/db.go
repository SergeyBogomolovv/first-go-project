package internal

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDB(cfg DBConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
