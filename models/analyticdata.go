package models

import "time"

type AnalyticData struct {
	AnalyticID     int       `json:"analytic_id"`
	ProductID      int       `json:"product_id"`
	ProductName    string    `json:"product_name,omitempty"`
	Quantity       int       `json:"quantity"`
	Revenue        float64   `json:"revenue"`
	Profit         float64   `json:"profit"`
	Contribution   float64   `json:"contribution"`
	IsTop20        bool      `json:"is_top_20"`
	AnalyticResult string    `json:"analytic_result"`
	AnalyticTime   time.Time `json:"analytic_time"`
}
