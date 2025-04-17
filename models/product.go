package models

type Product struct {
	ProductID   uint    `gorm:"primaryKey" json:"product_id"`
	ProductName string  `gorm:"not null" json:"product_name"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Note        string  `json:"note"`
}
