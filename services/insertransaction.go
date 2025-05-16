package services

import (
	"errors"
	"project_minyak/models"
	"time"

	"gorm.io/gorm"
)

func CreateTransactionWithDetailsBulk(db *gorm.DB, userID uint, fullname string, carts []models.Cart) (*models.Transaction, error) {
	if len(carts) == 0 {
		return nil, errors.New("keranjang kosong")
	}

	transaction := models.Transaction{
		UserID:            userID,
		ProductID:         carts[0].ProductID,
		UserFullname:      fullname,
		ProductName:       carts[0].Product.ProductName,
		StatusTransaction: "Pending",
	}
	if err := db.Create(&transaction).Error; err != nil {
		return nil, err
	}

	for _, item := range carts {
		detail := models.TransactionDetail{
			TransactionID: transaction.TransactionID,
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			Price:         item.Product.Price,
			DateTime:      time.Now(),
		}
		if err := db.Create(&detail).Error; err != nil {
			return nil, err
		}
	}

	return &transaction, nil
}
