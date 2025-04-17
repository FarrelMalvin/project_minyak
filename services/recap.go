package services

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"project_minyak/models"
)

func SalesRecap(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")

	query := `
		SELECT 
			DATE(td.date_time) as sale_date,
			td.Product_ID,
			p.product_name,
			SUM(td.quantity) as quantity_sold,
			p.Price,
			SUM(td.quantity) * p.Price as total_sales
		FROM transaction_details td
		JOIN transaction t ON td.Transaction_ID = t.Transaction_ID
		JOIN product p ON td.Product_ID = p.Product_ID
	`

	var args []interface{}
	if startDate != "" && endDate != "" {
		query += " WHERE DATE(td.date_time) BETWEEN $1 AND $2"
		args = append(args, startDate, endDate)
	}

	query += " GROUP BY sale_date, td.Product_ID, p.product_name, p.Price ORDER BY sale_date ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recaps []models.Recap
	for rows.Next() {
		var r models.Recap
		if err := rows.Scan(&r.SaleDate, &r.ProductID, &r.ProductName, &r.Quantity, &r.Price, &r.TotalSales); err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		recaps = append(recaps, r)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recaps)
}
