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
	Host string
	Port uint16
	User string
	Password string
	DB string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime string
}

func InitConfig() *Config {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	return &Config{
		Addr: env.GetString("ADDR", ":3000"),
		DB: DBConfig{
			Host: env.GetString("POSTGRES_HOST", "localhost"),
			Port: uint16(env.GetInt("POSTGRES_PORT", 5432)),
			User: env.GetString("POSTGRES_USER", "admin"),
			Password: env.GetString("POSTGRES_PASSWORD", "admin"),
			DB: env.GetString("POSTGRES_DB", "social"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONS", 30),
			MaxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}