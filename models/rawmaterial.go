package models

type RawMaterial struct {
	RawMaterialID   uint    `gorm:"primaryKey"`
	RawMaterialName string  `gorm:"not null"`
	Price           float64 `gorm:"type:decimal(10,2);not null"`
	Supplier        string  `gorm:"not null"`
}
