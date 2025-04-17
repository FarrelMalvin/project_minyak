package models

type Invoice struct {
	InvoiceID     uint    `gorm:"primaryKey"`
	TransactionID uint    `gorm:"not null"`
	TotalPrice    float64 `gorm:"type:decimal(10,2);not null"`
	PaymentMethod string  `gorm:"type:enum('Cash','Credit Card','Bank Transfer','E-Wallet');not null"`
}
