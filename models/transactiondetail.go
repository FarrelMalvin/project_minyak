package models

import "time"

type TransactionDetail struct {
	TransactionDetailID uint      `gorm:"primaryKey"`
	TransactionID       uint      `gorm:"not null;index"`
	ProductID           uint      `gorm:"not null;index"`
	Quantity            int       `gorm:"not null"`
	Price               float64   `gorm:"type:decimal(10,2);not null"`
	DateTime            time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Product     Product     `gorm:"foreignKey:ProductID;references:ProductID"`
	Transaction Transaction `gorm:"foreignKey:TransactionID" json:"-"`
}

func (TransactionDetail) TableName() string {
	return "transaction_detail"
}
