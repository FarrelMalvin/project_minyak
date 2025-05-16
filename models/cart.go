package models

import "time"

type Cart struct {
	CartID    uint      `gorm:"primaryKey;autoIncrement" json:"cart_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	AddedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"added_at"`

	// Relasi opsional (jika ingin preload)
	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (Cart) TableName() string {
	return "cart"
}
