package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"notif/internal/calendar/graph"
)

var studentMu sync.Mutex

const studentFilePath = "studentreq.json"

// loadStudentRequests reads all requests from studentreq.json
func loadStudentRequests() ([]graph.StudentRequest, error) {
	data, err := os.ReadFile(studentFilePath)
	if os.IsNotExist(err) {
		return []graph.StudentRequest{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []graph.StudentRequest{}, nil
	}
	var requests []graph.StudentRequest
	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, err
	}
	return requests, nil
}

// saveStudentRequests writes all requests back to studentreq.json
func saveStudentRequests(requests []graph.StudentRequest) error {
	data, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(studentFilePath, data, 0644)
}

// AppendStudentRequest saves a new student visit request and returns its ID
func AppendStudentRequest(req graph.StudentRequest) (string, error) {
	studentMu.Lock()
	defer studentMu.Unlock()

	requests, err := loadStudentRequests()
	if err != nil {
		return "", err
	}

	req.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	req.SubmittedAt = time.Now().Format("2006-01-02 15:04:05")

	requests = append(requests, req)
	return req.ID, saveStudentRequests(requests)
}

// GetRequestsByTA returns all student requests for a specific TA's office hours
// TAs call this to see who is coming and why
func GetRequestsByTA(taUID string) ([]graph.StudentRequest, error) {
	studentMu.Lock()
	defer studentMu.Unlock()

	requests, err := loadStudentRequests()
	if err != nil {
		return nil, err
	}

	// Cross-reference with officehours.json to find slots belonging to this TA
	officeHours, err := load() // reuses the load() from storage.go
	if err != nil {
		return nil, err
	}

	// Build a set of Outlook event IDs that belong to this TA
	taEventIDs := make(map[string]bool)
	for _, slot := range officeHours {
		if slot.TAUID == taUID {
			taEventIDs[slot.OutlookEventID] = true
		}
	}

	var result []graph.StudentRequest
	for _, r := range requests {
		if taEventIDs[r.OutlookEventID] {
			result = append(result, r)
		}
	}
	return result, nil
}

// GetRequestsByEvent returns all student requests for a specific calendar event
func GetRequestsByEvent(outlookEventID string) ([]graph.StudentRequest, error) {
	studentMu.Lock()
	defer studentMu.Unlock()

	requests, err := loadStudentRequests()
	if err != nil {
		return nil, err
	}

	var result []graph.StudentRequest
	for _, r := range requests {
		if r.OutlookEventID == outlookEventID {
			result = append(result, r)
		}
	}
	return result, nil
}

// GetRequestsByStudent returns all requests submitted by a specific student
func GetRequestsByStudent(studentUID string) ([]graph.StudentRequest, error) {
	studentMu.Lock()
	defer studentMu.Unlock()

	requests, err := loadStudentRequests()
	if err != nil {
		return nil, err
	}

	var result []graph.StudentRequest
	for _, r := range requests {
		if r.StudentUID == studentUID {
			result = append(result, r)
		}
	}
	return result, nil
}

// DeleteStudentRequest removes a request by its local ID
// Useful if a student wants to cancel their visit
func DeleteStudentRequest(id string) error {
	studentMu.Lock()
	defer studentMu.Unlock()

	requests, err := loadStudentRequests()
	if err != nil {
		return err
	}

	filtered := requests[:0]
	for _, r := range requests {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}

	if len(filtered) == len(requests) {
		return fmt.Errorf("request with id %s not found", id)
	}
	return saveStudentRequests(filtered)
}
