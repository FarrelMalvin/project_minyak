package services

import (
	"encoding/json"
	"net/http"

	"project_minyak/config"
	"project_minyak/models"

	"gorm.io/gorm"
)

func AddToCart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, _ := config.ExtractClaimsFromRequest(r)
		userID := claims.UserID

		var req struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		cart := models.Cart{
			UserID:    uint(userID),
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := db.Create(&cart).Error; err != nil {
			http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cart)
	}
}

func GetUserCart(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, _ := config.ExtractClaimsFromRequest(r)
		userID := claims.UserID

		var cart []models.Cart
		if err := db.Preload("Product").Where("user_id = ?", userID).Find(&cart).Error; err != nil {
			http.Error(w, "Failed to get cart", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cart)
	}
}

func DeleteCartItems(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			CartIDs []uint `json:"cart_ids"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		if err := db.Where("cart_id IN ?", req.CartIDs).Delete(&models.Cart{}).Error; err != nil {
			http.Error(w, "Failed to delete cart items", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
