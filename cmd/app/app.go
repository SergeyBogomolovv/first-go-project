package app

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	DB *sql.DB
	Router *chi.Mux
}

type Module interface {
	Register(db *sql.DB, router *chi.Mux)
}

func (app *Application) RegisterModule(module Module) {
	module.Register(app.DB, app.Router)
}

func (app *Application) Run(addr string) error {
	s := &http.Server{
		Addr: addr,
		Handler: app.Router,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 15,
		IdleTimeout: time.Minute,
	}

	log.Printf("Server start on %s", addr)

	return s.ListenAndServe()
}

