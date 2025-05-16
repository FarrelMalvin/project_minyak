package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"project_minyak/config"
	"project_minyak/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

func CheckoutHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := config.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		userID := uint(claims.UserID)
		customerName := claims.Name
		customerEmail := claims.Email

		var req struct {
			Items []struct {
				ProductID uint    `json:"product_id"`
				Quantity  int     `json:"quantity"`
				Price     float64 `json:"price"`
				Name      string  `json:"name"`
			} `json:"items"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(req.Items) == 0 {
			http.Error(w, "No items provided", http.StatusBadRequest)
			return
		}

		var carts []models.Cart
		for _, item := range req.Items {
			// Validasi manual untuk item kosong
			if item.ProductID == 0 || item.Quantity <= 0 || item.Price <= 0 || item.Name == "" {
				http.Error(w, "Invalid item data", http.StatusBadRequest)
				return
			}
			carts = append(carts, models.Cart{
				UserID:    userID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Product: models.Product{
					ProductID:   item.ProductID,
					ProductName: item.Name,
					Price:       float64(item.Price),
				},
			})
		}

		// Simpan transaksi
		transaction, err := CreateTransactionWithDetailsBulk(db, userID, customerName, carts)
		if err != nil {
			log.Println("❌ Gagal insert transaksi:", err)
			http.Error(w, "Gagal menyimpan transaksi", http.StatusInternalServerError)
			return
		}

		orderID := fmt.Sprintf("ORDER-%d", transaction.TransactionID)
		var snapClient snap.Client
		snapClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

		var totalAmount int64
		var items []midtrans.ItemDetails
		for _, item := range carts {
			itemTotal := int64(item.Product.Price) * int64(item.Quantity)
			totalAmount += itemTotal

			items = append(items, midtrans.ItemDetails{
				ID:    fmt.Sprintf("%d", item.ProductID),
				Price: int64(item.Product.Price),
				Qty:   int32(item.Quantity),
				Name:  item.Product.ProductName,
			})
		}

		snapReq := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: totalAmount,
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: customerName,
				Email: customerEmail,
			},
			Items: &items,
		}

		snapResp, err := snapClient.CreateTransaction(snapReq)
		if err != nil {
			log.Println("❌ Midtrans error:", err)
			http.Error(w, "Midtrans error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":      "Transaksi berhasil dibuat",
			"token":        snapResp.Token,
			"redirect_url": snapResp.RedirectURL,
		})
	}
}
