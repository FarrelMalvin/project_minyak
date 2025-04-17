package models

type Recap struct {
	SaleDate    string  `json:"sale_date"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity_sold"`
	Price       float64 `json:"price"`
	TotalSales  float64 `json:"total_sales"`
}
