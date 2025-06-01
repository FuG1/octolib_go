package BookHandlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"octolib/api/handlers/BooksHandlers"
	"octolib/api/models"
	"octolib/api/services"
	"octolib/db"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := db.InitDB(); err != nil {
		panic(err)
	}
	// Очистка таблиц
	db.DB.Exec("TRUNCATE TABLE books RESTART IDENTITY CASCADE")
	db.DB.Exec("TRUNCATE TABLE authors RESTART IDENTITY CASCADE")
	db.DB.Exec("TRUNCATE TABLE genres RESTART IDENTITY CASCADE")

	// Добавление тестовых данных
	db.DB.Exec("INSERT INTO authors (first_name, last_name) VALUES ('John', 'Doe')")
	db.DB.Exec("INSERT INTO genres (name) VALUES ('Fiction')")

	services.JwtKey = []byte("testSecretKey")
	m.Run()
}

func createTestJWT(role int) string {
	claims := &services.Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(services.JwtKey)
	return signed
}

func TestAddBook(t *testing.T) {
	book := models.Book{
		Title:         "Test Book",
		Author:        1,
		Genre:         1, // Убедитесь, что genre_id существует
		PublishedDate: "2023-05-31",
	}
	data, _ := json.Marshal(book)
	req, _ := http.NewRequest(http.MethodPost, "/api/books/add", bytes.NewBuffer(data))
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: createTestJWT(2)})

	rr := httptest.NewRecorder()
	BookHandlers.AddBookHandler(rr, req)
	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestUpdateBook(t *testing.T) {
	updatedBook := models.Book{
		Title:         "Updated Title",
		Author:        1,
		Genre:         1, // Убедитесь, что genre_id существует
		PublishedDate: "2023-06-01",
	}
	data, _ := json.Marshal(updatedBook)
	req, _ := http.NewRequest(http.MethodPut, "/api/books/update?id=1", bytes.NewBuffer(data))
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: createTestJWT(2)})

	rr := httptest.NewRecorder()
	BookHandlers.UpdateBookHandler(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteBook(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "/api/books/delete?id=1", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: createTestJWT(2)})

	rr := httptest.NewRecorder()
	BookHandlers.DeleteBookHandler(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}
