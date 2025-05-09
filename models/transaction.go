package models

type Transaction struct {
	TransactionID     uint    `gorm:"primaryKey"`
	UserID            uint    `gorm:"not null"`
	ProductID         uint    `gorm:"not null"`
	UserFullname      string  `gorm:"not null"`
	ProductName       string  `gorm:"not null"`
	StatusTransaction string  `gorm:"type:enum('Pending','Completed','Cancelled');not null"`
	Product           Product `gorm:"foreignKey:ProductID" json:"product"`
	User              User    `gorm:"foreignKey:UserID" json:"user"`
}

func (Transaction) TableName() string {
	return "transaction"
}
