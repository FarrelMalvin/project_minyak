package services

import (
	"encoding/json"
	"fmt"
	"log"
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

		// Ambil token dari Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token dan ambil data user
		claims, err := config.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		userID := uint(claims.UserID)
		customerName := claims.Name
		customerEmail := claims.Email

		// Ambil data dari body request
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
		if req.ProductID == 0 || req.Quantity <= 0 || req.Price <= 0 || req.ProductName == "" {
			http.Error(w, "Missing or invalid data", http.StatusBadRequest)
			return
		}

		// Hitung total
		grossAmount := int64(req.Price * float64(req.Quantity))

		// Simpan transaksi ke database
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
			log.Println("❌ Gagal insert transaksi:", err)
			http.Error(w, "Database error saat menyimpan transaksi", http.StatusInternalServerError)
			return
		}
		orderID := fmt.Sprintf("ORDER-%d", transaction.TransactionID)

		// Ambil server key dari environment
		serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
		if serverKey == "" {
			http.Error(w, "Midtrans server key not found", http.StatusInternalServerError)
			return
		}

		// Inisialisasi Midtrans Snap client
		var snapClient snap.Client
		snapClient.New(serverKey, midtrans.Sandbox) // Ganti ke midtrans.Production untuk production

		// Buat data transaksi untuk Midtrans
		snapReq := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: grossAmount,
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: customerName,
				Email: customerEmail,
			},
			Items: &[]midtrans.ItemDetails{
				{
					ID:    fmt.Sprintf("%d", req.ProductID),
					Price: int64(req.Price),
					Qty:   int32(req.Quantity),
					Name:  req.ProductName,
				},
			},
		}

		// Buat transaksi Midtrans
		snapResp, err := snapClient.CreateTransaction(snapReq)
		if err != nil {
			log.Printf("❌ Midtrans error: %v\n", err)
			http.Error(w, "Gagal membuat transaksi Midtrans", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"message":      "Transaksi berhasil dibuat",
			"token":        snapResp.Token,
			"redirect_url": snapResp.RedirectURL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("❌ Gagal encode response:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

	}
}
