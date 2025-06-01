package AuthHandlers

import (
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"octolib/db"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := db.InitDB()
	if err != nil {
		panic(err)
	}

	_, err = db.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		panic(err)
	}

	_, err = db.DB.Exec(`
        CREATE TABLE IF NOT EXISTS users(
            id SERIAL PRIMARY KEY,
            username VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            role_id INT NOT NULL
        )
    `)
	if err != nil {
		panic(err)
	}

	// Добавляем тестового пользователя с хешированным паролем
	hashed, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	_, err = db.DB.Exec(`
        INSERT INTO users(username, password, role_id)
        VALUES($1, $2, $3)
        ON CONFLICT (username) DO NOTHING
    `, "testuser", string(hashed), 1)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	db.DB.Close()
	os.Exit(code)
}

func TestLoginHandler_Success(t *testing.T) {
	body := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
}

func TestLoginHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/login", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestRegisterHandler_Success(t *testing.T) {
	body := map[string]string{
		"username": "newuser",
		"password": "newpassword123",
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
	}
}

func TestRegisterHandler_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/register", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestRegisterHandler_MissingFields(t *testing.T) {
	body := map[string]string{
		"username": "",
		"password": "",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestRegisterHandler_ShortPassword(t *testing.T) {
	body := map[string]string{
		"username": "testuser2",
		"password": "short",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestRegisterHandler_InvalidUsername(t *testing.T) {
	body := map[string]string{
		"username": "invalid*user",
		"password": "validpass123",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestRegisterHandler_InvalidPassword(t *testing.T) {
	body := map[string]string{
		"username": "validuser123",
		"password": "pas!@#",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestRegisterHandler_UserAlreadyExists(t *testing.T) {
	// Хэшируем пароль вручную, чтобы не создать конфликт в базе
	hashed, _ := bcrypt.GenerateFromPassword([]byte("newpassword1234"), bcrypt.DefaultCost)
	_, _ = db.DB.Exec(`INSERT INTO users (username, password, role_id) VALUES ($1, $2, 1)`, "testuser2", string(hashed))

	body := map[string]string{
		"username": "testuser2",
		"password": "newpassword1234",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, status)
	}
}
