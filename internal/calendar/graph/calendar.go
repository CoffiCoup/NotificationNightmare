package graph

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type CalendarEvent struct {
	TAName    string
	Day       string // "2026-04-15"
	StartTime string // "14:00"
	EndTime   string // "15:00"
	Location  string
}

// CreateCalendarEvent creates an Outlook event and returns the Outlook event ID
func CreateCalendarEvent(event CalendarEvent) (string, error) {
	token, err := GetCalendarToken()
	if err != nil {
		return "", fmt.Errorf("failed to get calendar token: %w", err)
	}

	calendarID := os.Getenv("OUTLOOK_CALENDAR_ID")
	endpoint := fmt.Sprintf("/me/calendars/%s/events", calendarID)

	payload := map[string]interface{}{
		"subject": fmt.Sprintf("Office Hours — %s", event.TAName),
		"start": map[string]string{
			"dateTime": fmt.Sprintf("%sT%s:00", event.Day, event.StartTime),
			"timeZone": "America/New_York",
		},
		"end": map[string]string{
			"dateTime": fmt.Sprintf("%sT%s:00", event.Day, event.EndTime),
			"timeZone": "America/New_York",
		},
		"location": map[string]string{
			"displayName": event.Location,
		},
		"showAs": "free",
	}

	resp, err := GraphRequest("POST", endpoint, token, payload)
	if err != nil {
		return "", fmt.Errorf("calendar request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return "", fmt.Errorf("calendar event creation failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse event ID: %w", err)
	}
	return result.ID, nil
}

// UpdateCalendarEvent modifies an existing Outlook event by its Outlook event ID
func UpdateCalendarEvent(outlookEventID string, event CalendarEvent) error {
	token, err := GetCalendarToken()
	if err != nil {
		return fmt.Errorf("failed to get calendar token: %w", err)
	}

	endpoint := fmt.Sprintf("/me/events/%s", outlookEventID)

	payload := map[string]interface{}{
		"subject": fmt.Sprintf("Office Hours — %s", event.TAName),
		"start": map[string]string{
			"dateTime": fmt.Sprintf("%sT%s:00", event.Day, event.StartTime),
			"timeZone": "America/New_York",
		},
		"end": map[string]string{
			"dateTime": fmt.Sprintf("%sT%s:00", event.Day, event.EndTime),
			"timeZone": "America/New_York",
		},
		"location": map[string]string{
			"displayName": event.Location,
		},
	}

	resp, err := GraphRequest("PATCH", endpoint, token, payload)
	if err != nil {
		return fmt.Errorf("update request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update failed with status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// DeleteCalendarEvent removes an Outlook event by its Outlook event ID
func DeleteCalendarEvent(outlookEventID string) error {
	token, err := GetCalendarToken()
	if err != nil {
		return fmt.Errorf("failed to get calendar token: %w", err)
	}

	endpoint := fmt.Sprintf("/me/events/%s", outlookEventID)

	resp, err := GraphRequest("DELETE", endpoint, token, nil)
	if err != nil {
		return fmt.Errorf("delete request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed with status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
