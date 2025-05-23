package config

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("secret_key")

type UserClaims struct {
	UserID uint   `json:"user_id"` // ubah dari int ke uint
	Role   string `json:"role"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(role string, userID uint, name string, email string) (string, error) {
	claims := UserClaims{
		Role:   role,
		UserID: userID,
		Name:   name,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func ExtractClaimsFromRequest(r *http.Request) (*UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("authorization header missing or invalid")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := ParseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	return claims, nil
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
