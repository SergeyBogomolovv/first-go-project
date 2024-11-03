package main

import (
	"context"
	"go-back/cmd/app"
	"go-back/internal/config"
	"go-back/internal/db"
	"go-back/internal/posts"
	"go-back/internal/users"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() { 
	cfg := config.InitConfig()
	db, err := db.ConnectToDB(cfg.DB)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("database connection established")
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	app := app.NewApplication(db, r) 

	usersModule := users.NewUsersModule()
	postsModule := posts.NewPostsModule()

	app.RegisterModule(usersModule)

	postsModule.InjectUserService(usersModule.Service)
	app.RegisterModule(postsModule)

	ctx, cancel := context.WithCancel(context.Background())
	
  defer cancel()

  go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	if err := app.Run(ctx, cfg.Addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}