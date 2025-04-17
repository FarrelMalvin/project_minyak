package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"project_minyak/models"
)

func InsertProductAndStock(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Product models.Product `json:"product"`
		Stock   models.Stock   `json:"stock"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate stock value
	if requestData.Stock.Stock <= 0 {
		http.Error(w, "Stock quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}

	// Insert product into the database
	productQuery := `INSERT INTO "product" (Product_name, price, note) VALUES ($1, $2, $3)`
	result, err := tx.Exec(productQuery, requestData.Product.ProductName, requestData.Product.Price, requestData.Product.Note)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to insert product: %v", err)
		http.Error(w, "Failed to insert product", http.StatusInternalServerError)
		return
	}

	productID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to get product_id: %v", err)
		http.Error(w, "Failed to get product_id", http.StatusInternalServerError)
		return
	}

	requestData.Stock.ProductID = uint(productID)

	// Insert stock into the database
	stockQuery := `INSERT INTO "stock" (Product_ID, Stock) VALUES ($1, $2)`
	_, err = tx.Exec(stockQuery, requestData.Stock.ProductID, requestData.Stock.Stock)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to insert stock: %v", err)
		http.Error(w, "Failed to insert stock", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Return a more detailed response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Product and stock inserted successfully",
		"product_id":     productID,
		"stock_quantity": requestData.Stock.Stock,
	})
}

func UpdateProductStock(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if stock.ProductID == 0 {
		http.Error(w, "ProductID is required", http.StatusBadRequest)
		return
	}

	query := `UPDATE stock SET stock = $1 WHERE product_id = $2`

	_, err = db.Exec(query, stock.Stock, stock.ProductID)
	if err != nil {
		http.Error(w, "Failed to update stock", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stock updated successfully"))
}
