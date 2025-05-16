package BookHandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"octolib/api/models"
	"octolib/db"
)

func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем, передан ли author_id
	if book.Author == 0 {
		http.Error(w, "Author ID is required", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли author_id в таблице authors
	row := db.DB.QueryRow("SELECT COUNT(*) FROM authors WHERE id = $1", book.Author)
	var count int
	if err := row.Scan(&count); err != nil || count == 0 {
		http.Error(w, "Author ID does not exist", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли книга с таким же названием от того же автора
	row = db.DB.QueryRow("SELECT COUNT(*) FROM books WHERE title = $1 AND author_id = $2", book.Title, book.Author)
	if err := row.Scan(&count); err != nil {
		log.Printf("Error checking book existence: %v", err)
		http.Error(w, "Error checking book existence", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Book with the same title by this author already exists", http.StatusConflict)
		return
	}

	// Добавляем книгу в таблицу books
	_, err := db.DB.Exec(
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
