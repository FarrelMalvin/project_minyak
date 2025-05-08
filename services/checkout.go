package services

import (
	"fmt"
	"project_minyak/models"
	"time"

	"gorm.io/gorm"
)

// Insert transaction & detail before calling Midtrans
func CreateTransactionWithDetails(db *gorm.DB, userID uint, productID uint, fullname, productName string, quantity int, price float64) (*models.Transaction, error) {
	fmt.Println("📝 Membuat transaksi baru...")
	fmt.Printf("UserID: %d, ProductID: %d, Name: %s, Product: %s, Qty: %d, Price: %.2f\n",
		userID, productID, fullname, productName, quantity, price)

	// Buat transaksi utama
	transaction := models.Transaction{
		UserID:            userID,
		ProductID:         productID,
		UserFullname:      fullname,
		ProductName:       productName,
		StatusTransaction: "Pending", // Status default sebelum pembayaran
	}

	// Simpan transaksi utama
	if err := db.Create(&transaction).Error; err != nil {
		fmt.Println("❌ Gagal menyimpan transaksi utama:", err)
		return nil, err
	}
	fmt.Println("✅ Transaksi utama disimpan. TransactionID:", transaction.TransactionID)

	// Buat detail transaksi
	detail := models.TransactionDetail{
		TransactionID: transaction.TransactionID,
		ProductID:     productID,
		Quantity:      quantity,
		Price:         price,
		DateTime:      time.Now(),
	}

	// Simpan detail transaksi
	if err := db.Create(&detail).Error; err != nil {
		fmt.Println("❌ Gagal menyimpan detail transaksi:", err)
		return nil, err
	}
	fmt.Println("✅ Detail transaksi disimpan.")

	return &transaction, nil
}
