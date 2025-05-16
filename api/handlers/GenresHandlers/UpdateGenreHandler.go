package GenresHandlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"octolib/api/models"
	"octolib/api/services"
	"octolib/db"
	"strconv"
)

func UpdateGenreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	genreIDStr := r.URL.Query().Get("id")
	if genreIDStr == "" {
		http.Error(w, "Genre ID is required", http.StatusBadRequest)
		return
	}

	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		http.Error(w, "Invalid Genre ID", http.StatusBadRequest)
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

	row := db.DB.QueryRow("SELECT COUNT(*) FROM genres WHERE id = $1", genreID)
	var count int
	if err := row.Scan(&count); err != nil || count == 0 {
		http.Error(w, "Genre not found", http.StatusNotFound)
		return
	}

	_, err = db.DB.Exec("UPDATE genres SET name = $1 WHERE id = $2", genre.Name, genreID)
	if err != nil {
		log.Printf("Error updating genre in database: %v", err)
		http.Error(w, "Error updating genre in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Genre updated successfully"))
}
