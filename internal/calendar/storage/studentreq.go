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

// GetRequestsByTA returns all student requests for a specific TA's slots
func GetRequestsByTA(taUID string) ([]graph.StudentRequest, error) {
	studentMu.Lock()
	defer studentMu.Unlock()

	requests, err := loadStudentRequests()
	if err != nil {
		return nil, err
	}

	// Get all slot IDs belonging to this TA
	mu.Lock()
	officeHours, err := load()
	mu.Unlock()
	if err != nil {
		return nil, err
	}

	taSlotIDs := make(map[string]bool)
	for _, slot := range officeHours {
		if slot.TAUID == taUID {
			taSlotIDs[slot.ID] = true
		}
	}

	var result []graph.StudentRequest
	for _, r := range requests {
		if taSlotIDs[r.SlotID] {
			result = append(result, r)
		}
	}
	return result, nil
}

// GetRequestsBySlot returns all student requests for a specific slot ID
func GetRequestsBySlot(slotID string) ([]graph.StudentRequest, error) {
	studentMu.Lock()
	defer studentMu.Unlock()

	requests, err := loadStudentRequests()
	if err != nil {
		return nil, err
	}

	var result []graph.StudentRequest
	for _, r := range requests {
		if r.SlotID == slotID {
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
