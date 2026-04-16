package models

// Profile: The static info about the TA
type Profile struct {
	ComputingID string `json:"computing_id"`
	Name        string `json:"name"`
	Title       string `json:"title"` // e.g., "Lead TA"
	Bio         string `json:"bio"`
	PhotoURL    string `json:"photo_url"`
	Email       string `json:"email"`
}

// Availability: The specific slots a TA is free (The Calendar data)
type Availability struct {
	ComputingID string `json:"computing_id"`
	Name        string `json:"name"`
	Date        string `json:"date"`
	StartTime   string `json:"start_time"`
	Duration    int    `json:"duration"` //in hours
	Location    string `json:"location"`
	Description string `json:"description"`
}

// OHRequest: When a student actually claims a slot
type OHRequest struct {
	ComputingID string `json:"computing_id"`
	TAID        string `json:"ta_id"`
	DateTime    string `json:"date_time"`
	Reason      string `json:"reason"`
}

//Page: Structure for management and security of webpages
type WebPage struct {
	Security int    //role # of user has to be <= to this to gain access to this page
	URL      string //what the http request needs to be directed to (relative filepath) TODO: NEEDS TESTING!
}
