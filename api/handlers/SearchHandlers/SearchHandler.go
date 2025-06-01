package SearchHandlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"octolib/db"
)

func SearchBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	code := query.Get("code")
	title := query.Get("title")

	if code == "" && title == "" {
		http.Error(w, "Either 'code' or 'title' parameter is required", http.StatusBadRequest)
		return
	}

	var row *sql.Row
	if code != "" {
		row = db.DB.QueryRow("SELECT title, author_id, genre_id, description, published_date, popularity, code FROM books WHERE code = $1", code)
	} else {
		row = db.DB.QueryRow("SELECT title, author_id, genre_id, description, published_date, popularity, code FROM books WHERE title = $1", title)
	}

	var book struct {
		Title         string `json:"title"`
		AuthorID      int    `json:"author_id"`
		GenreID       int    `json:"genre_id"`
		Description   string `json:"description"`
		PublishedDate string `json:"published_date"`
		Popularity    int    `json:"popularity"`
		Code          string `json:"code"`
	}
	if err := row.Scan(&book.Title, &book.AuthorID, &book.GenreID, &book.Description, &book.PublishedDate, &book.Popularity, &book.Code); err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
