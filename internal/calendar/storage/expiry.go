package storage

import (
	"fmt"
	"log"
	"time"
)

func StartExpiryWorker() {
	go func() {
		log.Println("Expiry worker started — checking every 30 seconds")
		for {
			if err := archiveExpiredSlots(); err != nil {
				log.Println("WARNING expiry worker error:", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()
}

func archiveExpiredSlots() error {
	slots, err := GetAllOfficeHours()
	if err != nil {
		return err
	}

	now := time.Now().In(time.FixedZone("America/New_York", -4*60*60))
	var expired []string

	for _, slot := range slots {
		endStr := fmt.Sprintf("%sT%s:00", slot.Day, slot.EndTime)
		endTime, err := time.ParseInLocation("2006-01-02T15:04:05",
			endStr, time.FixedZone("America/New_York", -4*60*60))
		if err != nil {
			log.Printf("WARNING could not parse end time for slot %s: %v", slot.ID, err)
			continue
		}

		log.Printf("DEBUG NOW: %v | END: %v", now, endTime)

		if now.After(endTime) {
			log.Printf("Archiving expired slot: %s %s %s-%s (SlotID: %s)",
				slot.TAName, slot.Day, slot.StartTime, slot.EndTime, slot.ID)

			if err := ArchiveRequestsForSlot(slot); err != nil {
				log.Printf("WARNING failed to archive slot %s: %v", slot.ID, err)
				continue
			}
			expired = append(expired, slot.ID)
		}
	}

	if len(expired) > 0 {
		if err := removeExpiredSlots(expired); err != nil {
			return err
		}
		log.Printf("Archived and removed %d expired slot(s)", len(expired))
	}

	return nil
}

func removeExpiredSlots(ids []string) error {
	mu.Lock()
	defer mu.Unlock()

	rows, err := load()
	if err != nil {
		return err
	}

	toRemove := make(map[string]bool)
	for _, id := range ids {
		toRemove[id] = true
	}

	filtered := rows[:0]
	for _, r := range rows {
		if !toRemove[r.ID] {
			filtered = append(filtered, r)
		}
	}

	return save(filtered)
}
