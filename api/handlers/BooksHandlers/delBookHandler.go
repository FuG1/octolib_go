package BookHandlers

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"octolib/api/services"
	"octolib/db"
	"strconv"
)

func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.DB.Exec("DELETE FROM books WHERE id = $1", bookID)
	if err != nil {
		log.Printf("Error deleting book from database: %v", err)
		http.Error(w, "Error deleting book from database", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec(`
		DO $$
		BEGIN
			CREATE TEMP TABLE temp_books AS SELECT * FROM books ORDER BY id;
			TRUNCATE books;
			INSERT INTO books SELECT row_number() OVER () AS id, title, author_id, genre_id, description, published_date, popularity FROM temp_books;
			DROP TABLE temp_books;
		END $$;
	`)
	if err != nil {
		log.Printf("Error reordering book IDs: %v", err)
		http.Error(w, "Error reordering book IDs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Book deleted"))
}
