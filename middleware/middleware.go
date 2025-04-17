package middleware

import (
	"context"
	"net/http"
	"strings"

	"project_minyak/config"
)

func extractAndVerifyToken(w http.ResponseWriter, r *http.Request, allowedRole string) (*config.UserClaims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return nil, false
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := config.ParseToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil, false
	}

	if claims.Role != allowedRole {
		http.Error(w, "Forbidden: Only "+allowedRole+" can access this resource", http.StatusForbidden)
		return nil, false
	}

	return claims, true
}

func CustomerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := extractAndVerifyToken(w, r, "Customer")
		if !ok {
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := extractAndVerifyToken(w, r, "Admin")
		if !ok {
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ManagerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := extractAndVerifyToken(w, r, "Manager")
		if !ok {
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SalesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := extractAndVerifyToken(w, r, "Sales")
		if !ok {
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
