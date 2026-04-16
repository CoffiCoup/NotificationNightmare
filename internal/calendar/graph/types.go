package graph

// OfficeHoursRow represents one office hours entry stored in officehours.json
type OfficeHoursRow struct {
	ID             string `json:"id"`
	TAUID          string `json:"ta_uid"`
	TAName         string `json:"ta_name"`
	Day            string `json:"day"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Location       string `json:"location"`
	OutlookEventID string `json:"outlook_event_id"`
}

// StudentRequest represents one student's visit request stored in studentreq.json
type StudentRequest struct {
	ID             string `json:"id"`               // unique local ID
	StudentUID     string `json:"student_uid"`      // from NetBadge session
	OutlookEventID string `json:"outlook_event_id"` // links to the specific office hours slot
	TAName         string `json:"ta_name"`
	Day            string `json:"day"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Location       string `json:"location"`
	Reason         string `json:"reason"` // student's reason for visiting
	SubmittedAt    string `json:"submitted_at"`
}
