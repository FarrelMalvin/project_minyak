package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_minyak/middleware"
	"project_minyak/models"
)

func ViewTransaction(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Ambil user ID dari context
	userIDInterface := r.Context().Value(middleware.UserIDKey)
	if userIDInterface == nil {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Cast ke int (kalau di klaim JWT kamu UserID-nya disimpan sebagai int)
	userID, ok := userIDInterface.(int)
	if !ok {
		log.Printf("DEBUG: userIDInterface type: %T\n", userIDInterface) // log tipe data
		http.Error(w, fmt.Sprintf("Invalid User ID format: %v", userIDInterface), http.StatusUnauthorized)
		return

	}
	log.Printf("DEBUG: userID: %d\n", userID)

	query := `
		SELECT 
			t.transaction_id, t.user_id, t.user_fullname, 
			t.product_id, p.product_name, t.status_transaction
		FROM transaction t
		INNER JOIN product p ON t.product_id = p.product_id
		WHERE t.user_id = $1
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error retrieving transactions: %v", err)
		http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(
			&transaction.TransactionID,
			&transaction.UserID,
			&transaction.UserFullname,
			&transaction.ProductID,
			&transaction.ProductName,
			&transaction.StatusTransaction,
		); err != nil {
			log.Printf("Error scanning transaction data: %v", err)
			http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
