package models

type Stock struct {
	StockID   uint    `gorm:"primaryKey" json:"stock_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Stock     int     `gorm:"not null" json:"stock"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (Stock) TableName() string {
	return "stock"
}

