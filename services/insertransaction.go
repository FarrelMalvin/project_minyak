package services

import (
	"project_minyak/models"
	"time"

	"gorm.io/gorm"
)

// Insert transaction & detail before calling Midtrans
func CreateTransactionWithDetails(db *gorm.DB, userID uint, productID uint, fullname, productName string, quantity int, price float64) (*models.Transaction, error) {
	// Buat transaksi utama
	transaction := models.Transaction{
		UserID:            userID,
		ProductID:         productID,
		UserFullname:      fullname,
		ProductName:       productName,
		StatusTransaction: "Pending", // Status default sebelum pembayaran
	}

	if err := db.Create(&transaction).Error; err != nil {
		return nil, err
	}

	// Buat detail transaksi
	detail := models.TransactionDetail{
		TransactionID: transaction.TransactionID,
		ProductID:     productID,
		Quantity:      quantity,
		Price:         price,
		DateTime:      time.Now(),
	}

	if err := db.Create(&detail).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}
