package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "go-back/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title			Golang API Example
// @version		1.0
// @description	This is my first Go API
// @contact.email	bogomolovs693@gmail.com
func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Get("/docs/*", httpSwagger.Handler())

	return r
}

type Application struct {
	Router *chi.Mux
	DB     *sqlx.DB
	Server *http.Server
}

type Module interface {
	Register(db *sqlx.DB, router *chi.Mux)
}

func (app *Application) RegisterModule(module Module) {
	module.Register(app.DB, app.Router)
}

func (app *Application) Run(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func() {
		log.Printf("Server started on %s", app.Server.Addr)
		errChan <- app.Server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.Server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		log.Println("Server gracefully stopped")
	case err := <-errChan:
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func NewApplication(addr string, db *sqlx.DB, router *chi.Mux) *Application {
	return &Application{DB: db, Router: router, Server: &http.Server{
		Addr:         addr,
		Handler:      router,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Minute,
	}}
}
