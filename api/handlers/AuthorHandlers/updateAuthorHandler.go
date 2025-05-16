package AuthorHandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"octolib/api/models"
	"octolib/api/services"
	"octolib/db"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
)

func UpdateAuthorHandler(w http.ResponseWriter, r *http.Request) {
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

	authorIDStr := r.URL.Query().Get("id")
	if authorIDStr == "" {
		http.Error(w, "Author ID is required", http.StatusBadRequest)
		return
	}

	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		http.Error(w, "Invalid Author ID", http.StatusBadRequest)
		return
	}

	var author models.Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if author.FirstName == "" || author.LastName == "" {
		http.Error(w, "First and Last name's can't be empty", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(
		"UPDATE authors SET first_name = $1, last_name = $2, middle_name = $3 WHERE id = $4",
		author.FirstName, author.LastName, author.MiddleName, authorID,
	)
	if err != nil {
		log.Printf("Error updating author in database: %v", err)
		http.Error(w, "Error updating author in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Author updated successfully"))
}
