package health

import (
	"encoding/json"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status":  "ok",
		"message": "Backend Toko mini berjalan",
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Gagal mengirim respons health: %v", err)
	}
}
