package AuthorHandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"octolib/api/models"
	"octolib/api/services"
	"octolib/db"

	"github.com/golang-jwt/jwt/v4"
)

func AddAuthorHandler(w http.ResponseWriter, r *http.Request) {
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

	row := db.DB.QueryRow(
		"SELECT COUNT(*) FROM authors WHERE first_name = $1 AND last_name = $2 AND middle_name = $3",
		author.FirstName, author.LastName, author.MiddleName,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		log.Printf("Error checking author existence: %v", err)
		http.Error(w, "Error checking author existence", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Author already exists", http.StatusConflict)
		return
	}

	// Добавление автора в базу данных
	_, err = db.DB.Exec(
		"INSERT INTO authors (first_name, last_name, middle_name) VALUES ($1, $2, $3)",
		author.FirstName, author.LastName, author.MiddleName,
	)
	if err != nil {
		log.Printf("Error saving author to database: %v", err)
		http.Error(w, "Error saving author to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Author added successfully"))
}
