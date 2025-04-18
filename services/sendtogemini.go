package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"project_minyak/models"
)

func SendToGeminiForAnalysis(recaps []models.Recap) (string, error) {
	apiKey := "AIzaSyAYaIvg2Hs46gvYk7_FOtlGdLc7QM9rEVA"

	// Langkah 1: Format prompt dan data
	recapJSON, _ := json.Marshal(recaps)
	prompt := `Tolong lakukan analisis Pareto dari data penjualan berikut. 
Berikan hasil *hanya* dalam bentuk JSON array seperti contoh berikut, tanpa penjelasan tambahan, tanpa markdown, tanpa teks lain:

[
  {
    "Product ID": 1,
    "Product Name": "Contoh",
    "Total Quantity Sold": 10,
    "Total Revenue": 10000.0,
    "Profit": 2000.0,
    "Top 20%": true
  },
  ...
]`

	// Payload format baru sesuai dokumentasi Gemini
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt + "\n\nData:\n" + string(recapJSON)},
				},
			},
		},
	}

	reqBody, _ := json.Marshal(payload)

	// Langkah 2: Kirim request POST
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + apiKey
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", errors.New("Gemini API error: " + string(bodyBytes))
	}

	// Langkah 3: Ambil respons teks
	var geminiResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	// Parsing kandidat pertama
	candidates, ok := geminiResp["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", errors.New("No candidates returned by Gemini")
	}

	content := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	text := parts[0].(map[string]interface{})["text"].(string)

	return text, nil
}
