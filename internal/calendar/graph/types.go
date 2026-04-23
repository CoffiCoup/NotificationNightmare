package graph

// OfficeHoursRow represents one office hours entry stored in officehours.json
type OfficeHoursRow struct {
	ID        string `json:"id"`
	TAUID     string `json:"ta_uid"`
	TAName    string `json:"ta_name"`
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
}

// StudentRequest represents one student's visit request stored in studentreq.json
type StudentRequest struct {
	ID          string `json:"id"`
	StudentUID  string `json:"student_uid"`
	SlotID      string `json:"slot_id"`
	TAName      string `json:"ta_name"`
	Day         string `json:"day"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Location    string `json:"location"`
	Reason      string `json:"reason"`
	SubmittedAt string `json:"submitted_at"`
}

// StoredRequest represents a finalized student request archived after office hours end
type StoredRequest struct {
	ID          string `json:"id"`
	StudentUID  string `json:"student_uid"`
	SlotID      string `json:"slot_id"`
	TAName      string `json:"ta_name"`
	TAUID       string `json:"ta_uid"`
	Day         string `json:"day"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Location    string `json:"location"`
	Reason      string `json:"reason"`
	SubmittedAt string `json:"submitted_at"`
	ArchivedAt  string `json:"archived_at"`
}

// TABio represents a TA's profile stored in tabios.json
type TABio struct {
	ID        string `json:"id"`
	TAUID     string `json:"ta_uid"`
	Name      string `json:"name"`
	Title     string `json:"title"` // "Professor" or "TA"
	Bio       string `json:"bio"`
	PhotoURL  string `json:"photo_url"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
