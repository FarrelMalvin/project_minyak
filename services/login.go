package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"project_minyak/config"
	"project_minyak/models"

	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("secret_key")

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var user models.User
		query := `SELECT user_id, firstname, lastname, email, password, role FROM "user" WHERE email = $1`
		log.Println("Executing query with email:", creds.Email)
		err := db.QueryRow(query, creds.Email).Scan(
			&user.UserID, &user.Firstname, &user.Lastname,
			&user.Email, &user.Password, &user.Role,
		)
		if err != nil {
			log.Println("Query error:", err)
			http.Error(w, "Invalid email", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		tokenString, err := config.GenerateJWT(user.Role, int(user.UserID), user.Firstname+" "+user.Lastname, user.Email)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		log.Println("Login berhasil:", user.Email)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
			"email": user.Email,
			"name":  user.Firstname + " " + user.Lastname,
		})
	}
}
