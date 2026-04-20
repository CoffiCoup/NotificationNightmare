package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"notif/internal/calendar/graph"
	"notif/internal/calendar/storage"
)

type HoursHandler struct {
	DB *sql.DB // kept for future use, not required right now
}

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

	// Create Outlook calendar event — get back the Outlook event ID
	outlookEventID, err := graph.CreateCalendarEvent(graph.CalendarEvent{
		TAName:    req.TAName,
		Day:       req.Day,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Location:  req.Location,
	})
	if err != nil {
		log.Println("ERROR calendar write failed:", err)
		http.Error(w, "calendar error", http.StatusInternalServerError)
		return
	}

	// Save to JSON with the Outlook event ID
	localID, err := storage.AppendOfficeHoursJSON(graph.OfficeHoursRow{
		TAUID:          taUID,
		TAName:         req.TAName,
		Day:            req.Day,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Location:       req.Location,
		OutlookEventID: outlookEventID,
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

// UpdateOfficeHours handles PATCH /api/hours?id=xxx
func (h *HoursHandler) UpdateOfficeHours(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	var req SubmitHoursRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Look up the existing entry to get the Outlook event ID
	entry, err := storage.GetByID(id)
	if err != nil {
		http.Error(w, "entry not found", http.StatusNotFound)
		return
	}

	// Update Outlook calendar event
	if err := graph.UpdateCalendarEvent(entry.OutlookEventID, graph.CalendarEvent{
		TAName:    entry.TAName,
		Day:       req.Day,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Location:  req.Location,
	}); err != nil {
		log.Println("ERROR updating calendar:", err)
		http.Error(w, "calendar update error", http.StatusInternalServerError)
		return
	}

	// Update JSON file
	if err := storage.UpdateOfficeHoursJSON(id, req.Day, req.StartTime, req.EndTime, req.Location); err != nil {
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

	// Look up the existing entry to get the Outlook event ID
	_, err := storage.GetByID(id)
	if err != nil {
		http.Error(w, "entry not found", http.StatusNotFound)
		return
	}

	// Delete from JSON file
	if err := storage.DeleteOfficeHoursJSON(id); err != nil {
		log.Println("ERROR deleting from json:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
