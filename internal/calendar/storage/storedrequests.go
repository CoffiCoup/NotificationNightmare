package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"notif/internal/calendar/graph"
)

var storedMu sync.Mutex

const storedFilePath = "storedrequests.json"

func loadStoredRequests() ([]graph.StoredRequest, error) {
	data, err := os.ReadFile(storedFilePath)
	if os.IsNotExist(err) {
		return []graph.StoredRequest{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []graph.StoredRequest{}, nil
	}
	var requests []graph.StoredRequest
	if err := json.Unmarshal(data, &requests); err != nil {
		return nil, err
	}
	return requests, nil
}

func saveStoredRequests(requests []graph.StoredRequest) error {
	data, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(storedFilePath, data, 0644)
}

// ArchiveRequestsForSlot moves student requests for a completed slot into storedrequests.json
// Holds both locks for the entire operation to prevent partial writes
func ArchiveRequestsForSlot(slot graph.OfficeHoursRow) error {
	// Always acquire in the same order: studentMu → storedMu
	// This order must never be reversed anywhere in the codebase
	studentMu.Lock()
	defer studentMu.Unlock()
	storedMu.Lock()
	defer storedMu.Unlock()

	// Load all student requests
	allRequests, err := loadStudentRequests()
	if err != nil {
		return fmt.Errorf("failed to load student requests: %w", err)
	}

	// Split into matching (for this slot) and remaining (everything else)
	var matching []graph.StudentRequest
	var remaining []graph.StudentRequest
	for _, r := range allRequests {
		if r.SlotID == slot.ID {
			matching = append(matching, r)
		} else {
			remaining = append(remaining, r)
		}
	}

	log.Printf("DEBUG: Found %d matching request(s) for slot %s", len(matching), slot.ID)

	// Always clean up studentreq.json even if there are no requests to archive
	if len(matching) == 0 {
		log.Printf("No student requests to archive for slot %s — slot had no sign-ups", slot.ID)
		return nil
	}

	// Load existing archived requests
	stored, err := loadStoredRequests()
	if err != nil {
		return fmt.Errorf("failed to load storedrequests: %w", err)
	}

	// Convert StudentRequest → StoredRequest
	now := time.Now().Format("2006-01-02 15:04:05")
	for _, r := range matching {
		stored = append(stored, graph.StoredRequest{
			ID:          r.ID,
			StudentUID:  r.StudentUID,
			SlotID:      slot.ID,
			TAName:      slot.TAName,
			TAUID:       slot.TAUID,
			Day:         slot.Day,
			StartTime:   slot.StartTime,
			EndTime:     slot.EndTime,
			Location:    slot.Location,
			Reason:      r.Reason,
			SubmittedAt: r.SubmittedAt,
			ArchivedAt:  now,
		})
	}

	// Write to storedrequests.json first
	if err := saveStoredRequests(stored); err != nil {
		return fmt.Errorf("failed to save storedrequests: %w", err)
	}

	// Only remove from studentreq.json after successful archive write
	if err := saveStudentRequests(remaining); err != nil {
		return fmt.Errorf("failed to update studentreq after archiving: %w", err)
	}

	log.Printf("Successfully archived %d request(s) for slot %s", len(matching), slot.ID)
	return nil
}

// DiscardRequestsForSlot removes all student requests for a slot without archiving
// Call this when a slot is deleted or modified
func DiscardRequestsForSlot(slotID string) error {
	studentMu.Lock()
	defer studentMu.Unlock()

	allRequests, err := loadStudentRequests()
	if err != nil {
		return err
	}

	remaining := allRequests[:0]
	discarded := 0
	for _, r := range allRequests {
		if r.SlotID != slotID {
			remaining = append(remaining, r)
		} else {
			discarded++
		}
	}

	if discarded > 0 {
		log.Printf("Discarding %d request(s) for deleted/modified slot %s", discarded, slotID)
	}
	return saveStudentRequests(remaining)
}

// GetStoredRequestsByTA returns all archived requests for a specific TA
func GetStoredRequestsByTA(taUID string) ([]graph.StoredRequest, error) {
	storedMu.Lock()
	defer storedMu.Unlock()

	stored, err := loadStoredRequests()
	if err != nil {
		return nil, err
	}

	var result []graph.StoredRequest
	for _, r := range stored {
		if r.TAUID == taUID {
			result = append(result, r)
		}
	}
	return result, nil
}

// GetStoredRequestsBySlot returns all archived requests for a specific slot
func GetStoredRequestsBySlot(slotID string) ([]graph.StoredRequest, error) {
	storedMu.Lock()
	defer storedMu.Unlock()

	stored, err := loadStoredRequests()
	if err != nil {
		return nil, err
	}

	var result []graph.StoredRequest
	for _, r := range stored {
		if r.SlotID == slotID {
			result = append(result, r)
		}
	}
	return result, nil
}
