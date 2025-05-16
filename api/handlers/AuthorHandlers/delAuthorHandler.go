package AuthorHandlers

import (
	"log"
	"net/http"
	"octolib/api/services"
	"octolib/db"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

func DelAuthorHandler(w http.ResponseWriter, r *http.Request) {
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

	authorIDStr := r.URL.Query().Get("id")
	log.Printf("Received id: %s", authorIDStr)
	if authorIDStr == "" {
		log.Println("Author ID is missing in the query parameters")
		http.Error(w, "Author ID is required", http.StatusBadRequest)
		return
	}

	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		http.Error(w, "Invalid Author ID", http.StatusBadRequest)
		return
	}

	// Удаление автора из базы данных
	_, err = db.DB.Exec("DELETE FROM authors WHERE id = $1", authorID)
	if err != nil {
		log.Printf("Error deleting author from database: %v", err)
		http.Error(w, "Error deleting author from database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Author deleted successfully"))
}
