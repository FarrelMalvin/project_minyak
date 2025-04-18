package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"project_minyak/models"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// SignUp Handler for user registration
func CreateAccount(db *sql.DB) http.HandlerFunc {
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

		// Hash password before storing
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword
		c := cases.Title(language.English)
		user.Role = c.String(strings.ToLower(user.Role))

		// Validate role input
		if user.Role != "Admin" && user.Role != "Manager" && user.Role != "Sales" {
			http.Error(w, "Invalid role", http.StatusBadRequest)
			return
		}

		// Insert to database
		query := `INSERT INTO "user" (email, firstname, lastname, username, password, role) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = db.Exec(query, user.Email, user.Firstname, user.Lastname, user.Username, user.Password, user.Role)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User registered successfully with role: %s", user.Role)
	}
}
