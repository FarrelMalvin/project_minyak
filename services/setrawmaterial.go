package services

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"project_minyak/models"
)

func InsertRawMaterial(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var rawMaterial models.RawMaterial

	err := json.NewDecoder(r.Body).Decode(&rawMaterial)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	query := `INSERT INTO raw_material (raw_material_name, price, supplier) VALUES ($1, $2, $3)`
	_, err = tx.Exec(query, rawMaterial.RawMaterialName, rawMaterial.Price, rawMaterial.Supplier)

	if err != nil {
		tx.Rollback()
		log.Printf("Failed to insert raw_material: %v", err)
		http.Error(w, "Failed to insert raw_material", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Raw Material inserted successfully"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetRawMaterialsSorted(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	query := `SELECT raw_material_name, price, supplier FROM "raw_material" ORDER BY raw_material_name ASC, price ASC`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to retrieve raw materials: %v", err)
		http.Error(w, "Failed to retrieve raw materials", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rawMaterials []models.RawMaterial

	for rows.Next() {
		var rawMaterial models.RawMaterial
		if err := rows.Scan(&rawMaterial.RawMaterialName, &rawMaterial.Price, &rawMaterial.Supplier); err != nil {
			log.Printf("Failed to scan raw material: %v", err)
			http.Error(w, "Failed to retrieve raw materials", http.StatusInternalServerError)
			return
		}
		rawMaterials = append(rawMaterials, rawMaterial)
	}

	// Mengembalikan data dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rawMaterials)
}
