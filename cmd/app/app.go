package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Application struct {
	DB *sqlx.DB
	Router *chi.Mux
}

type Module interface {
	Register(db *sqlx.DB, router *chi.Mux)
}

func (app *Application) RegisterModule(module Module) {
	module.Register(app.DB, app.Router)
}

func (app *Application) Run(ctx context.Context, addr string) error {
	s := &http.Server{
		Addr: addr,
		Handler: app.Router,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 15,
		IdleTimeout: time.Minute,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		s.Shutdown(context.Background())
	}()

	log.Printf("Server start on %s", addr)

	return s.ListenAndServe()
}

func NewApplication(db *sqlx.DB, router *chi.Mux) *Application {
	return &Application{DB: db, Router: router}
}