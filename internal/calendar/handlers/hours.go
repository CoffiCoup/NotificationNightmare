package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"notif/internal/calendar/graph"
	"notif/internal/calendar/storage"
)

type HoursHandler struct{}

type SubmitHoursRequest struct {
	TAName    string `json:"ta_name"`
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
}

// SubmitOfficeHours handles POST /api/hours
func (h *HoursHandler) SubmitOfficeHours(w http.ResponseWriter, r *http.Request) {
	var req SubmitHoursRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.TAName == "" || req.Day == "" || req.StartTime == "" || req.EndTime == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// TODO: replace with real TA UID from NetBadge session
	taUID := "uid-placeholder"

	// Save directly to JSON — no Outlook call needed
	localID, err := storage.AppendOfficeHoursJSON(graph.OfficeHoursRow{
		TAUID:     taUID,
		TAName:    req.TAName,
		Day:       req.Day,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Location:  req.Location,
	})
	if err != nil {
		log.Println("ERROR json write failed:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     localID,
	})
}

// GetMyHours handles GET /api/hours?ta_uid=xxx
func (h *HoursHandler) GetMyHours(w http.ResponseWriter, r *http.Request) {
	taUID := r.URL.Query().Get("ta_uid")
	if taUID == "" {
		taUID = "uid-placeholder"
	}

	entries, err := storage.GetOfficeHoursByTA(taUID)
	if err != nil {
		log.Println("ERROR fetching hours:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// GetAllHours handles GET /api/hours/all
// Used by the HTML calendar template to fetch and render all office hours
func (h *HoursHandler) GetAllHours(w http.ResponseWriter, r *http.Request) {
	entries, err := storage.GetAllOfficeHours()
	if err != nil {
		log.Println("ERROR fetching all hours:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// UpdateOfficeHours handles PATCH /api/hours?id=xxx
func (h *HoursHandler) UpdateOfficeHours(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	// Decode the full replacement struct from the request body
	var updated graph.OfficeHoursRow
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if updated.TAName == "" || updated.Day == "" || updated.StartTime == "" || updated.EndTime == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// Discard all student requests for this slot — details may have changed
	if err := storage.DiscardRequestsForSlot(id); err != nil {
		log.Println("WARNING failed to discard student requests on update:", err)
	}

	// Replace the entire entry
	if err := storage.UpdateOfficeHoursJSON(id, updated); err != nil {
		log.Println("ERROR updating json:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// DeleteOfficeHours handles DELETE /api/hours?id=xxx
func (h *HoursHandler) DeleteOfficeHours(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	// Discard all student requests for this slot — TA is canceling
	if err := storage.DiscardRequestsForSlot(id); err != nil {
		log.Println("WARNING failed to discard student requests:", err)
		// non-fatal — continue with deletion
	}

	// Delete from JSON file
	if err := storage.DeleteOfficeHoursJSON(id); err != nil {
		log.Println("ERROR deleting from json:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
