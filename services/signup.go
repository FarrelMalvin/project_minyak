package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"project_minyak/models"

	"golang.org/x/crypto/bcrypt"
)

// Hash Password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func SignUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword

		user.Role = "Customer"

		query := `INSERT INTO "user" (firstname, lastname, email, username, password, role) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = db.Exec(query, user.Firstname, user.Lastname, user.Email, user.Username, user.Password, user.Role)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User registered successfully with role: Customer")
	}
}
