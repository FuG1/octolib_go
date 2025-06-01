package AuthorHandlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"octolib/api/handlers/AuthorHandlers"
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
	// Очищаем таблицу authors
	_, _ = db.DB.Exec("TRUNCATE TABLE authors RESTART IDENTITY CASCADE")
	// Устанавливаем тестовый ключ
	services.JwtKey = []byte("testSecretKey")

	// Запуск тестов
	m.Run()
}

func createJwtWithRole(role int) string {
	claims := &services.Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(services.JwtKey)
	return tokenStr
}

func TestAddAuthorHandler_Success(t *testing.T) {
	author := models.Author{FirstName: "John", LastName: "Doe"}
	body, _ := json.Marshal(author)
	req, err := http.NewRequest(http.MethodPost, "/api/author/add", bytes.NewBuffer(body))
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: createJwtWithRole(3)})
	rr := httptest.NewRecorder()

	AuthorHandlers.AddAuthorHandler(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", rr.Code)
	}
}

func TestDelAuthorHandler_Success(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/api/author/del?id=1", nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: createJwtWithRole(3)})
	rr := httptest.NewRecorder()

	AuthorHandlers.DelAuthorHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rr.Code)
	}
}

func TestUpdateAuthorHandler_Success(t *testing.T) {
	updated := models.Author{FirstName: "Jane", LastName: "Smith"}
	body, _ := json.Marshal(updated)
	req, err := http.NewRequest(http.MethodPut, "/api/author/update?id=1", bytes.NewBuffer(body))
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: createJwtWithRole(3)})
	rr := httptest.NewRecorder()

	AuthorHandlers.UpdateAuthorHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rr.Code)
	}
}
