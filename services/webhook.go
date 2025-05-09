package services

import (
	"encoding/json"
	"log"
	"net/http"
	"project_minyak/models"
	"strings"

	"gorm.io/gorm"
)

type MidtransNotification struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
}

func MidtransWebhookHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload MidtransNotification
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Println("Failed to decode webhook payload:", err)
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		log.Println("Webhook received. OrderID:", payload.OrderID, "Status:", payload.TransactionStatus)

		orderID := strings.TrimPrefix(payload.OrderID, "ORDER-") // Optional if you add prefix
		var transaction models.Transaction
		if err := db.Where("transaction_id = ?", orderID).First(&transaction).Error; err != nil {
			log.Println("Transaction not found:", err)
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}

		// Update status jika berhasil
		if payload.TransactionStatus == "capture" || payload.TransactionStatus == "settlement" {
			transaction.StatusTransaction = "Completed"
			if err := db.Save(&transaction).Error; err != nil {
				log.Println("Failed to update transaction status:", err)
				http.Error(w, "Failed to update status", http.StatusInternalServerError)
				return
			}

			// Kurangi stok produk
			var details []models.TransactionDetail
			db.Where("transaction_id = ?", transaction.TransactionID).Find(&details)
			for _, detail := range details {
				db.Model(&models.Stock{}).
					Where("product_id = ?", detail.ProductID).
					Update("stock", gorm.Expr("stock - ?", detail.Quantity))
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
