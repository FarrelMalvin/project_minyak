package models

type ParetoResult struct {
	ParetoID      uint    `gorm:"primaryKey;column:pareto_id" json:"pareto_id"`
	BatchID       uint    `gorm:"column:batch_id" json:"batch_id"`
	ProductID     uint    `gorm:"column:product_id" json:"product_id"`
	ProductName   string  `gorm:"column:product_name" json:"product_name"`
	TotalQuantity int     `gorm:"column:total_quantity" json:"total_quantity"`
	TotalRevenue  float64 `gorm:"column:total_revenue" json:"total_revenue"`
	Profit        float64 `gorm:"column:profit" json:"profit"`
	Contribution  float64 `gorm:"column:contribution" json:"contribution"`
	IsTop20       bool    `gorm:"column:is_top_20" json:"is_top_20"`
}

func (ParetoResult) TableName() string {
	return "pareto_analysis_result"
}
