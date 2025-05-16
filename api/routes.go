package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"octolib/api/handlers"
	customMiddleware "octolib/api/middlewares"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Общие middleware
	r.Use(middleware.Logger)

	// Маршруты без AuthMiddleware
	r.Post("/api/login", handlers.LoginHandler)
	r.Post("/api/register", handlers.RegisterHandler)

	// Группа маршрутов с AuthMiddleware
	r.Group(func(r chi.Router) {
		r.Use(customMiddleware.AuthMiddleware)
		//r.Post("/add-book", handlers.AddBookHandler) // Пример защищённого маршрута
	})

	return r
}
