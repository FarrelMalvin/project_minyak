package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"project_minyak/models"

	"gorm.io/gorm"
)

type DeleteUserRequest struct {
	UserID int `json:"user_id"`
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Coba ambil userID dari query param
		userIDStr := r.URL.Query().Get("id")
		var userID int
		var err error

		if userIDStr != "" {
			userID, err = strconv.Atoi(userIDStr)
			if err != nil || userID <= 0 {
				http.Error(w, "Invalid user ID in query parameter", http.StatusBadRequest)
				return
			}
		} else {

			var req DeleteUserRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}
			userID = req.UserID
			if userID <= 0 {
				http.Error(w, "Invalid user ID in JSON body", http.StatusBadRequest)
				return
			}
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Role != "Manager" && user.Role != "Sales" {
			http.Error(w, "Only Manager or Sales can be deleted", http.StatusForbidden)
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User deleted successfully"))
	}
}
