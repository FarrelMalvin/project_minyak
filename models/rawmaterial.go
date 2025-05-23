package models

type RawMaterial struct {
	RawMaterialID   uint    `gorm:"primaryKey" json:"raw_material_id"`
	RawMaterialName string  `gorm:"not null" json:"raw_material_name"`
	Price           float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Supplier        string  `gorm:"not null" json:"supplier"`
}

func (RawMaterial) TableName() string {
	return "rawmaterial"
}
