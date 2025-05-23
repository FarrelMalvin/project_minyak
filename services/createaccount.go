package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project_minyak/models"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

func CreateAccount(db *gorm.DB) http.HandlerFunc {
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

		// Password hashing
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword

		// Format Role: lower to Title case
		c := cases.Title(language.English)
		user.Role = c.String(strings.ToLower(user.Role))

		if user.Role != "Admin" && user.Role != "Manager" && user.Role != "Sales" {
			http.Error(w, "Invalid role", http.StatusBadRequest)
			return
		}

		// Insert using GORM
		result := db.Create(&user)
		if result.Error != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User registered successfully with role: %s", user.Role)
	}
}
