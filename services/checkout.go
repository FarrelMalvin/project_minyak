package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"project_minyak/config"
	"project_minyak/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

func CheckoutHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("Panic recovered: %v\nStack: %s", rec, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Ambil token dari header Authorization
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

		// Validasi item & buat list cart
		var carts []models.Cart
		for _, item := range req.Items {
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
					Price:       item.Price,
				},
			})
		}

		// Simpan transaksi ke DB
		transaction, err := CreateTransactionWithDetailsBulk(db, userID, customerName, carts)
		if err != nil {
			log.Println("Gagal insert transaksi:", err)
			http.Error(w, "Gagal menyimpan transaksi", http.StatusInternalServerError)
			return
		}

		// Midtrans
		log.Printf("transaction object: %+v\n", transaction)

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
		if snapResp == nil {
			log.Println("snapResp nil")
			http.Error(w, "Midtrans response nil", http.StatusInternalServerError)
			return
		}
		if snapResp.Token == "" || snapResp.RedirectURL == "" {
			log.Printf("snapResp invalid: %+v\n", snapResp)
			http.Error(w, "Gagal membuat transaksi Midtrans", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"message":      "Transaksi berhasil dibuat",
			"token":        snapResp.Token,
			"redirect_url": snapResp.RedirectURL,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Gagal encode response:", err)
			http.Error(w, "Gagal mengirim response", http.StatusInternalServerError)
		}

	}
}
