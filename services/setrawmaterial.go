package services

import (
	"encoding/json"
	"log"
	"net/http"
	"project_minyak/models"

	"gorm.io/gorm"
)

func InsertRawMaterial(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rawMaterial models.RawMaterial
		if err := json.NewDecoder(r.Body).Decode(&rawMaterial); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if rawMaterial.RawMaterialName == "" || rawMaterial.Price <= 0 || rawMaterial.Supplier == "" {
			http.Error(w, "Missing or invalid raw material data", http.StatusBadRequest)
			return
		}

		if err := db.Create(&rawMaterial).Error; err != nil {
			log.Printf("Failed to insert raw material: %v", err)
			http.Error(w, "Failed to insert raw material", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Raw Material inserted successfully"})
	}
}

func GetRawMaterialsSorted(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rawMaterials []models.RawMaterial
		if err := db.Order("raw_material_name ASC, price ASC").Find(&rawMaterials).Error; err != nil {
			log.Printf("Failed to retrieve raw materials: %v", err)
			http.Error(w, "Failed to retrieve raw materials", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rawMaterials)
	}
}
