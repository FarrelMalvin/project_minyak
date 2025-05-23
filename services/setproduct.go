package services

import (
	"encoding/json"
	"net/http"
	"project_minyak/models"

	"gorm.io/gorm"
)

func InsertProductAndStock(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Product models.Product `json:"product"`
		Stock   models.Stock   `json:"stock"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.Stock.Stock <= 0 {
		http.Error(w, "Stock quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// Insert product
		if err := tx.Create(&requestData.Product).Error; err != nil {
			return err
		}

		// Assign foreign key
		requestData.Stock.ProductID = requestData.Product.ProductID

		// Insert stock
		if err := tx.Create(&requestData.Stock).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		http.Error(w, "Failed to insert product and stock", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Product and stock inserted successfully",
		"product_id":     requestData.Product.ProductID,
		"stock_quantity": requestData.Stock.Stock,
	})
}

func UpdateProductStock(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var stock models.Stock
	if err := json.NewDecoder(r.Body).Decode(&stock); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if stock.ProductID == 0 {
		http.Error(w, "ProductID is required", http.StatusBadRequest)
		return
	}

	if err := db.Model(&models.Stock{}).
		Where("product_id = ?", stock.ProductID).
		Update("stock", stock.Stock).Error; err != nil {
		http.Error(w, "Failed to update stock", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stock updated successfully"))
}
