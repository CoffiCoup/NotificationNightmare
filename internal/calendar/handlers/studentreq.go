package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"notif/internal/calendar/graph"
	"notif/internal/calendar/storage"
)

type StudentReqHandler struct{}

type SubmitRequestBody struct {
	SlotID    string `json:"slot_id"`
	TAName    string `json:"ta_name"`
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
	Reason    string `json:"reason"`
}

// SubmitRequest handles POST /api/requests
// Called when a student submits a visit request
func (h *StudentReqHandler) SubmitRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body SubmitRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.SlotID == "" || body.Reason == "" {
		http.Error(w, "outlook_event_id and reason are required", http.StatusBadRequest)
		return
	}

	// TODO: replace with real student UID from NetBadge session
	studentUID := "student-uid-placeholder"

	id, err := storage.AppendStudentRequest(graph.StudentRequest{
		StudentUID: studentUID,
		SlotID:     body.SlotID,
		TAName:     body.TAName,
		Day:        body.Day,
		StartTime:  body.StartTime,
		EndTime:    body.EndTime,
		Location:   body.Location,
		Reason:     body.Reason,
	})
	if err != nil {
		log.Println("ERROR saving student request:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     id,
	})
}

// GetRequestsForTA handles GET /api/requests?ta_uid=xxx
// TAs call this to see all incoming student requests for their office hours
func (h *StudentReqHandler) GetRequestsForTA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taUID := r.URL.Query().Get("ta_uid")
	if taUID == "" {
		taUID = "uid-placeholder"
	}

	requests, err := storage.GetRequestsByTA(taUID)
	if err != nil {
		log.Println("ERROR fetching requests for TA:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// GetRequestsForSlot handles GET /api/requests?event_id=xxx
// Returns all student requests for a specific office hours slot
func (h *StudentReqHandler) GetRequestsForSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slotID := r.URL.Query().Get("slot_id")
	if slotID == "" {
		http.Error(w, "event_id is required", http.StatusBadRequest)
		return
	}

	requests, err := storage.GetRequestsBySlot(slotID)
	if err != nil {
		log.Println("ERROR fetching requests for event:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// GetMyRequests handles GET /api/requests/me?student_uid=xxx
// Students can see their own submitted requests
func (h *StudentReqHandler) GetMyRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: replace with real student UID from NetBadge session
	studentUID := r.URL.Query().Get("student_uid")
	if studentUID == "" {
		studentUID = "student-uid-placeholder"
	}

	requests, err := storage.GetRequestsByStudent(studentUID)
	if err != nil {
		log.Println("ERROR fetching student requests:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// DeleteRequest handles DELETE /api/requests?id=xxx
// Students can cancel their own visit request
func (h *StudentReqHandler) DeleteRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	if err := storage.DeleteStudentRequest(id); err != nil {
		log.Println("ERROR deleting student request:", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
