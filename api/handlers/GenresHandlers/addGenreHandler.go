package GenresHandlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"octolib/api/models"
	"octolib/api/services"
	"octolib/db"
)

func AddGenreHandler(w http.ResponseWriter, r *http.Request) {
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

	if claims.Role != 3 {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var genre models.Genre
	if err := json.NewDecoder(r.Body).Decode(&genre); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if genre.Name == "" {
		http.Error(w, "Genre name can't be empty", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли жанр с таким именем
	row := db.DB.QueryRow("SELECT COUNT(*) FROM genres WHERE name = $1", genre.Name)
	var count int
	if err := row.Scan(&count); err != nil {
		log.Printf("Error checking genre existence: %v", err)
		http.Error(w, "Error checking genre existence", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Genre already exists", http.StatusConflict)
		return
	}

	// Добавляем жанр в базу данных
	_, err = db.DB.Exec("INSERT INTO genres (name) VALUES ($1)", genre.Name)
	if err != nil {
		log.Printf("Error saving genre to database: %v", err)
		http.Error(w, "Error saving genre to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Genre added successfully"))
}
