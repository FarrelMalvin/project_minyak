package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project_minyak/middleware"
	"project_minyak/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type TransactionSummary struct {
	TransactionID uint      `json:"transaction_id"`
	Date          time.Time `json:"date"`
	TotalPrice    float64   `json:"total_price"`
	Status        string    `json:"status"`
	ItemsSummary  string    `json:"items_summary"`
}

func ViewTransactionSummary(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDInterface := r.Context().Value(middleware.UserIDKey)
		if userIDInterface == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID := userIDInterface.(uint)

		var transactions []models.Transaction
		if err := db.Preload("TransactionDetails.Product").
			Where("user_id = ?", userID).
			Find(&transactions).Error; err != nil {
			http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
			return
		}

		var response []TransactionSummary
		for _, trx := range transactions {
			var totalPrice float64
			itemMap := make(map[string]int)
			var date time.Time

			for _, detail := range trx.TransactionDetails {
				totalPrice += float64(detail.Quantity) * detail.Price
				itemMap[detail.Product.ProductName] += detail.Quantity
				if date.IsZero() {
					date = detail.DateTime
				}
			}

			var summaryItems string
			for name, qty := range itemMap {
				summaryItems += fmt.Sprintf("%dx %s, ", qty, name)
			}
			if len(summaryItems) > 2 {
				summaryItems = summaryItems[:len(summaryItems)-2]
			}

			response = append(response, TransactionSummary{
				TransactionID: trx.TransactionID,
				TotalPrice:    totalPrice,
				Status:        trx.StatusTransaction,
				ItemsSummary:  summaryItems,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func ViewTransactionDetailByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		transactionID, _ := strconv.Atoi(idStr)

		var details []models.TransactionDetail
		if err := db.Preload("Product").
			Where("transaction_id = ?", transactionID).
			Find(&details).Error; err != nil {
			http.Error(w, "Failed to retrieve details", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(details)
	}
}
