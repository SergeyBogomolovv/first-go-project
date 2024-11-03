package main

import (
	"go-back/cmd/app"
	"go-back/internal/config"
	"go-back/internal/db"
	"go-back/internal/users"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() { 
	cfg := config.InitConfig()
	db, err := db.ConnectToDB(cfg.DB)

	if err != nil {
		log.Panic(err)
	}

	log.Println("database connection established")
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	app := &app.Application{DB: db, Router: r}

	app.RegisterModule(users.NewUsersModule())

	log.Fatal(app.Run(cfg.Addr))
}