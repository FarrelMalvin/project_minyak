package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"project_minyak/models"
)

func GetStock(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	query := `
		SELECT 
			p.product_id, p.product_name, p.price, p.note, 
			s.stock
		FROM product p
		INNER JOIN stock s ON p.product_id = s.product_id
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error retrieving stock: %v", err)
		http.Error(w, "Failed to retrieve stock", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stocks []models.Stock

	for rows.Next() {
		var stock models.Stock

		if err := rows.Scan(
			&stock.Product.ProductID,
			&stock.Product.ProductName,
			&stock.Product.Price,
			&stock.Product.Note,
			&stock.Stock,
		); err != nil {
			log.Printf("Error scanning stock data: %v", err)
			http.Error(w, "Failed to scan stock data", http.StatusInternalServerError)
			return
		}

		stock.ProductID = stock.Product.ProductID
		stocks = append(stocks, stock)
	}

	// Mengatur response dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}
