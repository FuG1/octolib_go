package services

import (
	"octolib/db"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte(db.Jwt)

type Claims struct {
	ID   int `json:"id"`
	Role int `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(id int, role int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		ID:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtKey)
}
