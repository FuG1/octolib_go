package AuthHandlers

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"octolib/api/models"
	"octolib/db"
	"regexp"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверка на пустые поля
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password cannot be empty", http.StatusBadRequest)
		return
	}

	// Проверка длины пароля
	if len(user.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	// Проверка формата username
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !usernameRegex.MatchString(user.Username) {
		http.Error(w, "Username can only contain English letters and numbers", http.StatusBadRequest)
		return
	}

	// Проверка формата пароля
	passwordRegex := regexp.MustCompile(`^[a-zA-Z0-9._]+$`)
	if !passwordRegex.MatchString(user.Password) {
		http.Error(w, "Password can only contain English letters, numbers, dots, and underscores", http.StatusBadRequest)
		return
	}

	// Проверка на существование пользователя
	var existingUser string
	err := db.DB.QueryRow("SELECT username FROM users WHERE username = $1", user.Username).Scan(&existingUser)
	if err == nil {
		http.Error(w, "User already exists", http.StatusUnauthorized)
		return
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Сохранение пользователя в базе данных
	_, err = db.DB.Exec("INSERT INTO users (username, password, role_id) VALUES ($1, $2, $3)", user.Username, user.Password, 1) // 1 - ID роли 'user'
	if err != nil {
		http.Error(w, "Error saving user to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
