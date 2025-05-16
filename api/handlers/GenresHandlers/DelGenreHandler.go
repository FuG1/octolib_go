package GenresHandlers

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"octolib/api/services"
	"octolib/db"
	"strconv"
)

func DeleteGenreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	// Удаляем книги, связанные с жанром
	_, err = db.DB.Exec("DELETE FROM books WHERE genre_id = $1", genreID)
	if err != nil {
		log.Printf("Error deleting books from database: %v", err)
		http.Error(w, "Error deleting books from database", http.StatusInternalServerError)
		return
	}

	// Удаляем жанр
	_, err = db.DB.Exec("DELETE FROM genres WHERE id = $1", genreID)
	if err != nil {
		log.Printf("Error deleting genre from database: %v", err)
		http.Error(w, "Error deleting genre from database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Genre and associated books deleted successfully"))
}
