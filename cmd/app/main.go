package main

import (
	"context"
	"go-back/internal/config"
	"go-back/internal/database"
	"go-back/internal/modules/posts"
	"go-back/internal/modules/users"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	db, err := database.ConnectToDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("database connection established")
	defer db.Close()

	router := NewRouter()

	app := NewApplication(":3000", db, router)

	usersModule := users.NewUsersModule()
	app.RegisterModule(usersModule)

	postsModule := posts.NewPostsModule()
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

	if err := app.Run(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
