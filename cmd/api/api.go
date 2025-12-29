package main

import (
	"log"
	"minio-go-s3/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config   config
	pasteSvc service.PasteService
	gc       *GarbageCollector
}

type config struct {
	addr    string
	baseURL string
	db      dbConfig
	blob    blobConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type blobConfig struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	bucket          string
	useSSL          bool
}

// ---- application methods
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Get("/pastes/{id}", app.GetPasteHandler)
		r.Post("/pastes", app.CreatePasteHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Println("Server Running on", app.config.addr)
	return srv.ListenAndServe()
}
