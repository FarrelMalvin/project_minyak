package services

import (
	"encoding/json"
	"net/http"
	"project_minyak/models"

	"gorm.io/gorm"
)

type StockResponse struct {
	StockID     uint    `json:"stock_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Note        string  `json:"note"`
}

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

		var response []StockResponse
		for _, stock := range stocks {
			response = append(response, StockResponse{
				StockID:     stock.StockID,
				ProductID:   stock.ProductID,
				ProductName: stock.Product.ProductName,
				Price:       stock.Product.Price,
				Quantity:    stock.Stock,
				Note:        stock.Product.Note,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
