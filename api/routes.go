package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"octolib/api/handlers/AuthHandlers"
	"octolib/api/handlers/AuthorHandlers"
	"octolib/api/handlers/BooksHandlers"
	"octolib/api/handlers/GenresHandlers"
	"octolib/api/handlers/SearchHandlers"
	customMiddleware "octolib/api/middlewares"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Общие middleware
	r.Use(middleware.Logger)

	// Маршруты без AuthMiddleware
	r.Post("/api/login", AuthHandlers.LoginHandler)
	r.Post("/api/register", AuthHandlers.RegisterHandler)
	r.Get("/api/search", SearchHandlers.SearchBookHandler)
	// Группа маршрутов с AuthMiddleware
	r.Group(func(r chi.Router) {
		r.Use(customMiddleware.AuthMiddleware)
		//Books
		r.Post("/api/add_book", BookHandlers.AddBookHandler)
		r.Delete("/api/del_book", BookHandlers.DeleteBookHandler)
		r.Put("/api/update_book", BookHandlers.UpdateBookHandler)
		//Authors
		r.Post("/api/add_author", AuthorHandlers.AddAuthorHandler)
		r.Delete("/api/del_author", AuthorHandlers.DelAuthorHandler)
		r.Put("/api/update_author", AuthorHandlers.UpdateAuthorHandler)
		//Genres
		r.Post("/api/add_genre", GenresHandlers.AddGenreHandler)
		r.Delete("/api/del_genre", GenresHandlers.DeleteGenreHandler)
		r.Put("/api/update_genre", GenresHandlers.UpdateGenreHandler)

	})

	return r
}
