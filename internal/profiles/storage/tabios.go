package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"notif/internal/calendar/graph"
)

var bioMu sync.Mutex

const bioFilePath = "tabios.json"

func loadBios() ([]graph.TABio, error) {
	data, err := os.ReadFile(bioFilePath)
	if os.IsNotExist(err) {
		return []graph.TABio{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []graph.TABio{}, nil
	}
	var bios []graph.TABio
	if err := json.Unmarshal(data, &bios); err != nil {
		return nil, err
	}
	return bios, nil
}

func saveBios(bios []graph.TABio) error {
	data, err := json.MarshalIndent(bios, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(bioFilePath, data, 0644)
}

// UpsertBio creates or fully replaces a TA's bio
func UpsertBio(bio graph.TABio) (string, error) {
	bioMu.Lock()
	defer bioMu.Unlock()

	bios, err := loadBios()
	if err != nil {
		return "", err
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	for i, b := range bios {
		if b.TAUID == bio.TAUID {
			bio.ID = b.ID
			bio.CreatedAt = b.CreatedAt
			bio.UpdatedAt = now
			bios[i] = bio
			return bio.ID, saveBios(bios)
		}
	}

	bio.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	bio.CreatedAt = now
	bio.UpdatedAt = now
	bios = append(bios, bio)
	return bio.ID, saveBios(bios)
}

// GetAllBios returns every TA bio
func GetAllBios() ([]graph.TABio, error) {
	bioMu.Lock()
	defer bioMu.Unlock()
	return loadBios()
}

// GetBioByTA returns a single TA's bio by UID
func GetBioByTA(taUID string) (*graph.TABio, error) {
	bioMu.Lock()
	defer bioMu.Unlock()

	bios, err := loadBios()
	if err != nil {
		return nil, err
	}

	for _, b := range bios {
		if b.TAUID == taUID {
			return &b, nil
		}
	}
	return nil, fmt.Errorf("no bio found for TA %s", taUID)
}

// DeleteBio removes a TA's bio by UID
func DeleteBio(taUID string) error {
	bioMu.Lock()
	defer bioMu.Unlock()

	bios, err := loadBios()
	if err != nil {
		return err
	}

	filtered := bios[:0]
	for _, b := range bios {
		if b.TAUID != taUID {
			filtered = append(filtered, b)
		}
	}

	if len(filtered) == len(bios) {
		return fmt.Errorf("no bio found for TA %s", taUID)
	}
	return saveBios(filtered)
}

// GetBioByID returns a single bio by its ID field
func GetBioByID(id string) (*graph.TABio, error) {
	bioMu.Lock()
	defer bioMu.Unlock()

	bios, err := loadBios()
	if err != nil {
		return nil, err
	}

	for _, b := range bios {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, fmt.Errorf("no bio found with id %s", id)
}
