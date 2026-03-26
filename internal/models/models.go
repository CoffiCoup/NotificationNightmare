package models

// Profile: The static info about the TA (The "Bio" box in your chart)
type Profile struct {
	ComputingID string `json:"computing_id"`
	Name        string `json:"name"`
	Title       string `json:"title"` // e.g., "Lead TA"
	Bio         string `json:"bio"`
	PhotoURL    string `json:"photo_url"`
}

// Availability: The specific slots a TA is free (The "Calendar" data)
type Availability struct {
	ComputingID string `json:"computing_id"`
	Name        string `json:"name"` // Added for readability
	Date        string `json:"date"`
	StartTime   string `json:"start_time"`
	Duration    int    `json:"duration"` // in minutes
	Location    string `json:"location"`
}

// OHRequest: When a student actually claims a slot
type OHRequest struct {
	ComputingID string `json:"computing_id"`
	TAID        string `json:"ta_id"`
	DateTime    string `json:"date_time"`
	Reason      string `json:"reason"`
}
