package middleware

import (
	"context"
	"net/http"
	"octolib/api/services"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware проверяет JWT-токен
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлечение токена из заголовка Authorization
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization token is required", http.StatusUnauthorized)
			return
		}

		// Парсинг токена и проверка подписи
		claims := &services.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return services.JwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавление данных из токена в контекст
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
