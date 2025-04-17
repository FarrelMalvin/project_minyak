package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"project_minyak/models"
)

func UpdateTransactionStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var status models.MidtransStatus
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update transaction status
	updateQuery := `
		UPDATE transactions 
		SET status_transaction = ?
		WHERE transaction_id = ?
	`
	_, err = tx.Exec(updateQuery, status.Status, status.TransactionID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating transaction status: %v", err)
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	// If payment is successful, update stock
	if status.Status == "successpayment" {
		// Get quantity from transaction
		var productID int
		var quantity int
		queryTransaction := `
			SELECT product_id, quantity 
			FROM transactions 
			WHERE transaction_id = ?
		`
		err = tx.QueryRow(queryTransaction, status.TransactionID).Scan(&productID, &quantity)
		if err != nil {
			tx.Rollback()
			log.Printf("Error getting transaction details: %v", err)
			http.Error(w, "Failed to process transaction", http.StatusInternalServerError)
			return
		}

		// Update stock
		updateStockQuery := `
			UPDATE product 
			SET stock = stock - ? 
			WHERE product_id = ?
		`
		_, err = tx.Exec(updateStockQuery, quantity, productID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating stock: %v", err)
			http.Error(w, "Failed to update stock", http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Failed to complete transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaction status updated successfully",
	})
}
