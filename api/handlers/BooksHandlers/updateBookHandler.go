package BookHandlers

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

func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
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

	if claims.Role == 1 {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	bookIDStr := r.URL.Query().Get("id")
	if bookIDStr == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		http.Error(w, "Invalid Book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли книга с таким ID
	row := db.DB.QueryRow("SELECT COUNT(*) FROM books WHERE id = $1", bookID)
	var count int
	if err := row.Scan(&count); err != nil || count == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Обновляем данные книги
	_, err = db.DB.Exec(
		`UPDATE books 
		SET title = $1, author_id = $2, genre_id = $3, description = $4, published_date = $5, popularity = $6 
		WHERE id = $7`,
		book.Title, book.Author, book.Genre, book.Description, book.PublishedDate, book.Popularity, bookID,
	)
	if err != nil {
		log.Printf("Error updating book in database: %v", err)
		http.Error(w, "Error updating book in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book updated successfully"))
}
