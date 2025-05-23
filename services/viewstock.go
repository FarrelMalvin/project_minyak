package services

import (
	"encoding/json"
	"net/http"
	"project_minyak/models"

	"gorm.io/gorm"
)

func GetStock(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var stocks []models.Stock

		if err := db.Preload("Product").Find(&stocks).Error; err != nil {
			http.Error(w, "Failed to retrieve stock", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stocks)
	}
}
