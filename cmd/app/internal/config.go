package internal

import (
	"go-back/pkg"
	"log"

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

func InitConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Addr: pkg.GetString("ADDR", ":3000"),
		DB: DBConfig{
			Host:     pkg.GetString("POSTGRES_HOST", "localhost"),
			Port:     uint16(pkg.GetInt("POSTGRES_PORT", 5432)),
			User:     pkg.GetString("POSTGRES_USER", "admin"),
			Password: pkg.GetString("POSTGRES_PASSWORD", "admin"),
			DBName:   pkg.GetString("POSTGRES_DB", "social"),
		},
	}
}
