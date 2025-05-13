package handlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"octolib/api/models"
	"octolib/db"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials models.User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Поиск пользователя в базе данных
	var storedPassword string
	var role int
	err := db.DB.QueryRow("SELECT password, role FROM users WHERE username = $1", credentials.Username).Scan(&storedPassword, &role)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Генерация JWT-токена

	// Возвращаем токен
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
