package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_minyak/models"
	"project_minyak/utils"
	"sort"
	"strings"
)

func AnalyzeAllProductsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AnalyzeAllProductsHandler HIT 🚀")
		startDate := r.URL.Query().Get("start")
		endDate := r.URL.Query().Get("end")

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

		// Parse ke struct
		analytics, summary, recommendation := ParseGeminiResult(result)
		if analytics == nil {
			http.Error(w, "Failed to parse Gemini response", http.StatusInternalServerError)
			return
		}

		// Simpan hasil batch
		batchID, err := SaveAnalysisBatch(db, startDate, endDate, summary, recommendation)
		if err != nil {
			http.Error(w, "Failed to save batch summary: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Simpan ke pareto_analysis_result
		if err := SaveParetoResultToDB(db, batchID, analytics); err != nil {
			http.Error(w, "Failed to save pareto analysis result: "+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"summary":        summary,
			"recommendation": recommendation,
			"details":        analytics,
		})
	}
}

func SaveAnalysisBatch(db *sql.DB, start, end, summary, recommendation string) (int, error) {
	var batchID int
	err := db.QueryRow(`
		INSERT INTO pareto_analysis_batch (start_date, end_date, summary, recommendation, created_at)
		VALUES ($1, $2, $3, $4, NOW()) RETURNING batch_id`,
		start, end, summary, recommendation,
	).Scan(&batchID)
	return batchID, err
}

func SaveParetoResultToDB(db *sql.DB, batchID int, analytics []models.AnalyticData) error {
	var totalRevenue float64
	type temp struct {
		data    models.AnalyticData
		revenue float64
	}
	var tempData []temp

	for _, a := range analytics {
		totalRevenue += a.Revenue
		tempData = append(tempData, temp{data: a, revenue: a.Revenue})
	}

	sort.Slice(tempData, func(i, j int) bool {
		return tempData[i].revenue > tempData[j].revenue
	})

	for _, t := range tempData {
		contribution := (t.revenue / totalRevenue) * 100

		// GUNAKAN NILAI DARI GEMINI
		isTop := t.data.IsTop20

		_, err := db.Exec(`
			INSERT INTO pareto_analysis_result
			(batch_id, product_id, product_name, total_quantity, total_revenue, contribution, is_top_20)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			batchID, t.data.ProductID, t.data.ProductName, t.data.Quantity,
			t.data.Revenue, contribution, isTop,
		)
		if err != nil {
			log.Printf("Error inserting pareto result for product %d: %v", t.data.ProductID, err)
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

func ParseGeminiResult(result string) ([]models.AnalyticData, string, string) {
	result = cleanJSONResult(result)
	log.Println("Gemini RAW Response:", result)

	var raw []map[string]interface{}
	err := json.Unmarshal([]byte(result), &raw)
	if err != nil {
		log.Println("Failed to parse Gemini result as JSON:", err)
		return nil, "", ""
	}

	var data []models.AnalyticData
	var summaryList []string
	var recommendationList []string

	for _, item := range raw {
		prodID := utils.ToInt(item["Product ID"])
		prodName := utils.ToString(item["Product Name"])
		quantity := utils.ToInt(item["Total Quantity Sold"])
		revenue := utils.ToFloat(item["Total Revenue"])
		isTop := utils.ToBool(item["Top 20%"])

		s := utils.ToString(item["Summary"])
		r := utils.ToString(item["Recommendation"])

		if s != "" {
			summaryList = append(summaryList, fmt.Sprintf("%s: %s", prodName, s))
		}
		if r != "" {
			recommendationList = append(recommendationList, fmt.Sprintf("%s: %s", prodName, r))
		}

		resultStr, _ := json.Marshal(item)
		data = append(data, models.AnalyticData{
			ProductID:      prodID,
			ProductName:    prodName,
			Quantity:       quantity,
			Revenue:        revenue,
			IsTop20:        isTop,
			AnalyticResult: string(resultStr),
		})
	}

	summary := strings.Join(summaryList, "\n\n")
	recommendation := strings.Join(recommendationList, "\n\n")

	return data, summary, recommendation
}
