package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"notif/internal/calendar/storage"
)

type StoredRequestsHandler struct{}

// GetStoredRequestsForTA handles GET /api/stored?ta_uid=xxx
// TAs use this to analyze historical office hours attendance
func (h *StoredRequestsHandler) GetStoredRequestsForTA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taUID := r.URL.Query().Get("ta_uid")
	if taUID == "" {
		taUID = "uid-placeholder"
	}

	requests, err := storage.GetStoredRequestsByTA(taUID)
	if err != nil {
		log.Println("ERROR fetching stored requests:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// GetStoredRequestsForEvent handles GET /api/stored?event_id=xxx
// Returns archived requests for a specific past office hours slot
func (h *StoredRequestsHandler) GetStoredRequestsForSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slotID := r.URL.Query().Get("slot_id")
	if slotID == "" {
		http.Error(w, "slot_id is required", http.StatusBadRequest)
		return
	}

	requests, err := storage.GetStoredRequestsBySlot(slotID)
	if err != nil {
		log.Println("ERROR fetching stored requests by event:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}
