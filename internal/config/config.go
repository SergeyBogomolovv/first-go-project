package config

import (
	"go-back/internal/env"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr string
	DB DBConfig
}

type DBConfig struct {
	Addr string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime string
}

func InitConfig() Config {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	return Config{
		Addr: env.GetString("ADDR", ":3000"),
		DB: DBConfig{
			Addr: env.GetString("DB_ADDR", "postgres://admin:admin@localhost:5432/social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONS", 30),
			MaxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}