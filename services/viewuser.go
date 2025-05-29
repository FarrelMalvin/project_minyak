package services

import (
	"encoding/json"
	"net/http"
	"project_minyak/models"
	"strings"

	"gorm.io/gorm"
)

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func ViewUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var users []models.User
		if err := db.
			Where("role IN ?", []string{"Manager", "Sales"}).
			Select("firstname, lastname, email, role").
			Find(&users).Error; err != nil {
			http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
			return
		}

		var resp []UserResponse
		for _, u := range users {
			fullName := strings.TrimSpace(u.Firstname + " " + u.Lastname)
			resp = append(resp, UserResponse{
				Name:  fullName,
				Email: u.Email,
				Role:  u.Role,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
