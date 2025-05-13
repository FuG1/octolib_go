package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"octolib/api/handlers"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/api/login", handlers.LoginHandler)
	r.Post("/api/register", handlers.RegisterHandler)

	return r
}
