package GenresHandlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"octolib/api/handlers/GenresHandlers"
	"octolib/api/models"
	"octolib/api/services"
	"octolib/db"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Инициализация БД
	if err := db.InitDB(); err != nil {
		panic(err)
	}
	// Очистка таблицы genres
	db.DB.Exec("TRUNCATE TABLE genres RESTART IDENTITY CASCADE")
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

func TestAddGenre(t *testing.T) {
	genre := models.Genre{Name: "Fiction"}
	data, _ := json.Marshal(genre)
	req, _ := http.NewRequest(http.MethodPost, "/api/genres/add", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: createTestJWT(3)})

	rr := httptest.NewRecorder()
	GenresHandlers.AddGenreHandler(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestUpdateGenre(t *testing.T) {
	// Добавляем жанр для обновления
	db.DB.Exec("INSERT INTO genres (name) VALUES ('Old Genre')")

	updatedGenre := models.Genre{Name: "Updated Genre"}
	data, _ := json.Marshal(updatedGenre)
	req, _ := http.NewRequest(http.MethodPut, "/api/genres/update?id=1", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: createTestJWT(3)})

	rr := httptest.NewRecorder()
	GenresHandlers.UpdateGenreHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteGenre(t *testing.T) {
	// Добавляем жанр для удаления
	db.DB.Exec("INSERT INTO genres (name) VALUES ('Genre to Delete')")

	req, _ := http.NewRequest(http.MethodDelete, "/api/genres/delete?id=1", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: createTestJWT(3)})

	rr := httptest.NewRecorder()
	GenresHandlers.DeleteGenreHandler(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
