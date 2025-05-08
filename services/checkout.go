package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"project_minyak/config"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"

	"gorm.io/gorm"
)

func CheckoutHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ambil token JWT dari header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token
		claims, err := config.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		userID := uint(claims.UserID)
		customerName := claims.Name
		customerEmail := claims.Email

		// Decode data produk dari frontend
		var req struct {
			ProductID   uint    `json:"product_id"`
			ProductName string  `json:"product_name"`
			Quantity    int     `json:"quantity"`
			Price       float64 `json:"price"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.ProductID == 0 || req.Quantity <= 0 || req.Price <= 0 {
			http.Error(w, "Missing or invalid data", http.StatusBadRequest)
			return
		}

		// Hitung total harga
		grossAmount := int64(req.Price * float64(req.Quantity))

		// Simpan transaksi ke DB
		transaction, err := CreateTransactionWithDetails(
			db,
			userID,
			req.ProductID,
			customerName,
			req.ProductName,
			req.Quantity,
			req.Price,
		)
		if err != nil {
			http.Error(w, "Gagal menyimpan transaksi ke database", http.StatusInternalServerError)
			return
		}

		// Gunakan ID transaksi sebagai order_id Midtrans
		orderID := fmt.Sprintf("ORDER-%d", transaction.TransactionID)

		// Midtrans setup
		serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
		if serverKey == "" {
			http.Error(w, "Server key not found", http.StatusInternalServerError)
			return
		}
		snapClient := snap.Client{}
		snapClient.New(serverKey, midtrans.Sandbox)

		params := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: grossAmount,
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: customerName,
				Email: customerEmail,
			},
		}

		snapResp, err := snapClient.CreateTransaction(params)
		if err != nil {
			http.Error(w, "Failed to create Midtrans transaction: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token":        snapResp.Token,
			"redirect_url": snapResp.RedirectURL,
		})
	}
}
