package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_minyak/models"
	"sort"
	"time"
)

func AnalyzeAllProductsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AnalyzeAllProductsHandler HIT 🚀")
		var recaps []models.Recap
		if err := json.NewDecoder(r.Body).Decode(&recaps); err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		// Kirim ke Gemini
		result, err := SendToGeminiForAnalysis(recaps)
		if err != nil {
			http.Error(w, "Analysis failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse result Gemini jadi struct
		analytics := ParseGeminiResult(result)

		// Simpan ke DB
		if err := SaveParetoAnalysisToDB(db, analytics); err != nil {
			http.Error(w, "Failed to save analysis: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Kirim hasil balik ke frontend juga kalau perlu
		json.NewEncoder(w).Encode(analytics)
	}
}

func SaveParetoAnalysisToDB(db *sql.DB, analytics []models.AnalyticData) error {
	// Langkah 1: Hitung total kontribusi dari semua produk
	var totalRevenue float64
	type temp struct {
		data    models.AnalyticData
		revenue float64
	}
	var tempData []temp

	for _, a := range analytics {
		var revenue float64
		_, err := fmt.Sscanf(a.AnalyticResult, "Total Quantity Sold: %d, Total Revenue: %f, Profit: %f", new(int), &revenue, new(float64))
		if err != nil {
			log.Printf("Failed to parse analytic result for product %d: %v", a.ProductID, err)
			continue
		}
		totalRevenue += revenue
		tempData = append(tempData, temp{data: a, revenue: revenue})
	}

	// Langkah 2: Urutkan berdasarkan revenue tertinggi
	sort.Slice(tempData, func(i, j int) bool {
		return tempData[i].revenue > tempData[j].revenue
	})

	// Langkah 3: Hitung kontribusi dan tandai top 20%
	var cumulative float64
	for _, t := range tempData {
		contribution := (t.revenue / totalRevenue) * 100
		cumulative += contribution
		isTop := cumulative <= 80.0 // Pareto: 80/20 rule

		_, err := db.Exec(`
			INSERT INTO "analytic" (Product_ID, Analytic_result, Contribution, Is_Top_20, Analytic_time)
			VALUES ($1, $2, $3, $4, $5)`,
			t.data.ProductID, t.data.AnalyticResult, contribution, isTop, time.Now(),
		)
		if err != nil {
			log.Printf("Error inserting analytic data for product %d: %v", t.data.ProductID, err)
			return err
		}
	}
	return nil
}

func ParseGeminiResult(result string) []models.AnalyticData {
	var raw []map[string]interface{}
	err := json.Unmarshal([]byte(result), &raw)
	if err != nil {
		log.Println("Failed to parse Gemini result as JSON:", err)
		return nil
	}

	var data []models.AnalyticData
	for _, item := range raw {
		prodID := int(item["Product ID"].(float64))
		prodName := item["Product Name"].(string)
		quantity := int(item["Total Quantity Sold"].(float64))
		revenue := item["Total Revenue"].(float64)
		profit := item["Profit"].(float64)
		isTop := false
		if val, ok := item["Top 20%"]; ok {
			isTop = val.(bool)
		}

		resultStr, _ := json.Marshal(item) // optional string hasil mentah
		data = append(data, models.AnalyticData{
			ProductID:      prodID,
			ProductName:    prodName,
			Quantity:       quantity,
			Revenue:        revenue,
			Profit:         profit,
			IsTop20:        isTop,
			AnalyticResult: string(resultStr),
		})
	}
	return data
}
