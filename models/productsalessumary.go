package models

type ProductSalesSummary struct {
	ProductID    uint
	ProductName  string
	TotalSold    int
	Price        float64
	TotalRevenue float64
}
