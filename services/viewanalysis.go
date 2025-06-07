package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"project_minyak/models"
	"project_minyak/utils"

	"gorm.io/gorm"
)

// --- Response Struct ---
type ParetoResultResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Revenue     float64 `json:"revenue"`
	IsTop20     bool    `json:"is_top_20"`
}

type AnalysisBatchResponse struct {
	BatchID        uint                   `json:"batch_id"`
	BatchName      string                 `json:"batch_name"`
	Summary        string                 `json:"summary"`
	Recommendation string                 `json:"recommendation"`
	Results        []ParetoResultResponse `json:"results"`
}

// --- Handler Function ---
func GetParetoAnalysisHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		batchIDStr := r.URL.Query().Get("batch_id")
		if batchIDStr == "" {
			http.Error(w, "Missing batch_id parameter", http.StatusBadRequest)
			return
		}

		batchID, err := strconv.Atoi(batchIDStr)
		if err != nil {
			http.Error(w, "Invalid batch_id", http.StatusBadRequest)
			return
		}

		var batch models.ParetoBatch
		if err := db.First(&batch, "batch_id = ?", batchID).Error; err != nil {
			http.Error(w, "Batch not found", http.StatusNotFound)
			return
		}

		var resultsRaw []models.ParetoResult
		if err := db.Where("batch_id = ?", batchID).Find(&resultsRaw).Error; err != nil {
			http.Error(w, "Failed to fetch analysis results", http.StatusInternalServerError)
			return
		}

		var results []ParetoResultResponse
		for _, item := range resultsRaw {
			results = append(results, ParetoResultResponse{
				ProductID:   item.ProductID,
				ProductName: item.ProductName,
				Quantity:    item.TotalQuantity,
				Revenue:     item.TotalRevenue,
				IsTop20:     item.IsTop20,
			})
		}

		batchName := utils.FormatBatchName(batch.StartDate, batch.EndDate)

		response := AnalysisBatchResponse{
			BatchID:        batch.BatchID,
			BatchName:      batchName,
			Summary:        batch.Summary,
			Recommendation: batch.Recommendation,
			Results:        results,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
