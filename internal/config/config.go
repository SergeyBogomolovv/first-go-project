package config

import (
	"fmt"
	"go-back/pkg/env"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

type Config struct {
	Addr string `json:"addr"`
	DB   DBConfig
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	return &Config{
		Addr: env.GetString("ADDR", ":3000"),
		DB: DBConfig{
			Host:     env.GetString("POSTGRES_HOST", "localhost"),
			Port:     uint16(env.GetInt("POSTGRES_PORT", 5432)),
			User:     env.GetString("POSTGRES_USER", "admin"),
			Password: env.GetString("POSTGRES_PASSWORD", "admin"),
			DBName:   env.GetString("POSTGRES_DB", "social"),
		},
	}, nil
}
