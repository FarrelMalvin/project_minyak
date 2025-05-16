package services

import (
	"errors"
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> 582553e2e5f10efd64d96859bf4d11168377381e
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
<<<<<<< HEAD
		ProductID:         carts[0].ProductID,
		UserFullname:      fullname,
		ProductName:       carts[0].Product.ProductName,
=======
		UserFullname:      fullname,
>>>>>>> 582553e2e5f10efd64d96859bf4d11168377381e
		StatusTransaction: "Pending",
	}
	if err := db.Create(&transaction).Error; err != nil {
		return nil, err
	}

<<<<<<< HEAD
	for _, item := range carts {
=======
	var totalPrice float64 = 0

	for _, item := range carts {
		subtotal := float64(item.Quantity) * item.Product.Price
		totalPrice += subtotal

>>>>>>> 582553e2e5f10efd64d96859bf4d11168377381e
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

<<<<<<< HEAD
=======
	invoice := models.Invoice{
		TransactionID:   transaction.TransactionID,
		TotalPrice:      totalPrice,
		MidtransOrderID: fmt.Sprintf("ORDER-%d", transaction.TransactionID),
		PaymentMethod:   "", 
		
	}

	if err := db.Create(&invoice).Error; err != nil {
		return nil, err
	}

>>>>>>> 582553e2e5f10efd64d96859bf4d11168377381e
	return &transaction, nil
}
