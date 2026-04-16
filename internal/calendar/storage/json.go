package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"notif/internal/calendar/graph"
)

var mu sync.Mutex

const filePath = "officehours.json"

// load reads all rows from the JSON file
func load() ([]graph.OfficeHoursRow, error) {
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return []graph.OfficeHoursRow{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []graph.OfficeHoursRow{}, nil
	}
	var rows []graph.OfficeHoursRow
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// save writes all rows back to the JSON file
func save(rows []graph.OfficeHoursRow) error {
	data, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// AppendOfficeHoursJSON adds a new entry and returns its generated ID
func AppendOfficeHoursJSON(row graph.OfficeHoursRow) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	rows, err := load()
	if err != nil {
		return "", err
	}

	// Generate a simple unique ID from timestamp
	// Generates a string of base 10 integer based on the time using UnixNano to use as a unique ID for calendar event
	row.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	rows = append(rows, row)
	return row.ID, save(rows)
}

// GetAllOfficeHours returns every entry in the JSON file
func GetAllOfficeHours() ([]graph.OfficeHoursRow, error) {
	mu.Lock()
	defer mu.Unlock()
	return load()
}

// GetOfficeHoursByTA returns all entries for a specific TA
func GetOfficeHoursByTA(taUID string) ([]graph.OfficeHoursRow, error) {
	mu.Lock()
	defer mu.Unlock()

	rows, err := load()
	if err != nil {
		return nil, err
	}

	var result []graph.OfficeHoursRow
	for _, r := range rows {
		if r.TAUID == taUID {
			result = append(result, r)
		}
	}
	return result, nil
}

// GetByID finds a single entry by its local ID
func GetByID(id string) (*graph.OfficeHoursRow, error) {
	mu.Lock()
	defer mu.Unlock()

	rows, err := load()
	if err != nil {
		return nil, err
	}

	for _, r := range rows {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("entry with id %s not found", id)
}

// UpdateOfficeHoursJSON updates day, time, and location for an entry by ID
func UpdateOfficeHoursJSON(id, day, startTime, endTime, location string) error {
	mu.Lock()
	defer mu.Unlock()

	rows, err := load()
	if err != nil {
		return err
	}

	found := false
	for i, r := range rows {
		if r.ID == id {
			rows[i].Day = day
			rows[i].StartTime = startTime
			rows[i].EndTime = endTime
			rows[i].Location = location
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("entry with id %s not found", id)
	}
	return save(rows)
}

// DeleteOfficeHoursJSON removes an entry by ID
func DeleteOfficeHoursJSON(id string) error {
	mu.Lock()
	defer mu.Unlock()

	rows, err := load()
	if err != nil {
		return err
	}

	filtered := rows[:0]
	for _, r := range rows {
		if r.ID != id {
			filtered = append(filtered, r)
		}
	}

	if len(filtered) == len(rows) {
		return fmt.Errorf("entry with id %s not found", id)
	}
	return save(filtered)
}
