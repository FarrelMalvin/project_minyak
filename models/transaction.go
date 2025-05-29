package models

import "time"

type Transaction struct {
	TransactionID      uint                `gorm:"primaryKey"`
	UserID             uint                `gorm:"not null;index"`
	UserFullname       string              `gorm:"not null"`
	StatusTransaction  string              `gorm:"type:enum('Pending','Completed','Cancelled');not null;index"`
	User               User                `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt          time.Time           `gorm:"autoCreateTime"`
	TransactionDetails []TransactionDetail `gorm:"foreignKey:TransactionID"`
}

func (Transaction) TableName() string {
	return "transaction"
}
