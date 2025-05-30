package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_minyak/models"
	"sort"
	"strings"
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
		if analytics == nil {
			http.Error(w, "Failed to parse Gemini response", http.StatusInternalServerError)
			return
		}

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
	var totalRevenue float64
	type temp struct {
		data    models.AnalyticData
		revenue float64
	}
	var tempData []temp

	for _, a := range analytics {
		var revenue float64
		_, err := fmt.Sscanf(a.AnalyticResult, `{"Product ID":%d,"Product Name":"%s","Total Quantity Sold":%d,"Total Revenue":%f,"Profit":%f,"Top 20%%":%t}`,
			new(int), new(string), new(int), &revenue, new(float64), new(bool))
		if err != nil {
			log.Printf("Failed to parse analytic result for product %d: %v", a.ProductID, err)
			revenue = a.Revenue // fallback ke nilai langsung
		}
		totalRevenue += revenue
		tempData = append(tempData, temp{data: a, revenue: revenue})
	}

	sort.Slice(tempData, func(i, j int) bool {
		return tempData[i].revenue > tempData[j].revenue
	})

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

// Clean markdown/json formatting
func cleanJSONResult(result string) string {
	result = strings.TrimSpace(result)

	// Cek apakah ada blok kode markdown
	if strings.HasPrefix(result, "```") {
		// Cari blok antara tiga backtick
		start := strings.Index(result, "[")
		end := strings.LastIndex(result, "]") + 1
		if start >= 0 && end > start {
			return result[start:end]
		}
	}

	// Jika tidak dalam blok markdown tapi masih JSON array biasa
	if strings.HasPrefix(result, "[") && strings.HasSuffix(result, "]") {
		return result
	}

	// Coba ambil substring dari [ sampai ] jika JSON tersembunyi di tengah narasi
	start := strings.Index(result, "[")
	end := strings.LastIndex(result, "]") + 1
	if start >= 0 && end > start {
		return result[start:end]
	}

	return ""
}

func ParseGeminiResult(result string) []models.AnalyticData {
	result = cleanJSONResult(result)
	log.Println("Gemini RAW Response:", result)

	var raw []map[string]interface{}
	err := json.Unmarshal([]byte(result), &raw)
	if err != nil {
		log.Println("Failed to parse Gemini result as JSON:", err)
		return nil
	}

	var data []models.AnalyticData
	for _, item := range raw {
		prodID := toInt(item["Product ID"])
		prodName := toString(item["Product Name"])
		quantity := toInt(item["Total Quantity Sold"])
		revenue := toFloat(item["Total Revenue"])
		profit := toFloat(item["Profit"])
		isTop := toBool(item["Top 20%"])

		resultStr, _ := json.Marshal(item)
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

func toInt(v interface{}) int {
	if val, ok := v.(float64); ok {
		return int(val)
	}
	return 0
}

func toFloat(v interface{}) float64 {
	if val, ok := v.(float64); ok {
		return val
	}
	return 0
}

func toString(v interface{}) string {
	if val, ok := v.(string); ok {
		return val
	}
	return ""
}

func toBool(v interface{}) bool {
	if val, ok := v.(bool); ok {
		return val
	}
	return false
}
