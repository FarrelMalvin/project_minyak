package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("secret_key")

type UserClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	UserID   int    `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(username, role string, userID int) (string, error) {
	claims := UserClaims{
		Username: username,
		Role:     role,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	fmt.Printf("DEBUG: parsed claims: %+v\n", claims)
	return claims, nil
}
