package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"octolib/api/models"
	"octolib/api/services"
	"octolib/db"

	"github.com/golang-jwt/jwt/v4"
)

func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenStr := cookie.Value
	claims := &services.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return services.JwtKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if claims.Role == 1 {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == 0 {
		http.Error(w, "Title and Author are required fields", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO books (title, author_id, genre_id, description, published_date, popularity) VALUES ($1, $2, $3, $4, $5, $6)",
		book.Title, book.Author, book.Genre, book.Description, book.PublishedDate, book.Popularity,
	)
	if err != nil {
		log.Printf("Error saving book to database: %v", err)
		http.Error(w, "Error saving book to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book added successfully"))
}
