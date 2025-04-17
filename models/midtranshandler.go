package models

type MidtransStatus struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}
