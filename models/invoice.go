package models

type Invoice struct {
	InvoiceID       uint    `gorm:"primaryKey" json:"invoice_id"`
	TransactionID   uint    `json:"transaction_id"`
	TotalPrice      float64 `json:"total_price"`
	MidtransOrderID string  `json:"midtrans_order_id"`
	PaymentMethod   string  `json:"payment_method"`
}

func (Invoice) TableName() string {
	return "invoice"
}
