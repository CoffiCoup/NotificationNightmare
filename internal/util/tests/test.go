package tests

import (
	"html/template"
	"log"
	"net/http"

	profhandlers "notif/cmd/server/handling"
	calhandlers "notif/internal/calendar/handlers"
	calstore "notif/internal/calendar/storage"
)

func TestCentral() {

	// Load templates from profiles package
	tmpl, err := template.ParseFiles(
		"internal/pages/profiles.html",
		"internal/pages/ProfileUpload.html",
	)
	if err != nil {
		log.Fatal("FATAL: could not load templates:", err)
	}

	calstore.StartExpiryWorker()

	// Calendar handlers — note the calhandlers alias
	h := &calhandlers.HoursHandler{}
	http.HandleFunc("/api/hours", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.SubmitOfficeHours(w, r)
		case http.MethodGet:
			h.GetMyHours(w, r)
		case http.MethodPatch:
			h.UpdateOfficeHours(w, r)
		case http.MethodDelete:
			h.DeleteOfficeHours(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/hours/all", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetAllHours(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	sr := &calhandlers.StudentReqHandler{}
	http.HandleFunc("/api/requests", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			sr.SubmitRequest(w, r)
		case http.MethodGet:
			if r.URL.Query().Get("slot_id") != "" {
				sr.GetRequestsForSlot(w, r)
			} else {
				sr.GetRequestsForTA(w, r)
			}
		case http.MethodDelete:
			sr.DeleteRequest(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/requests/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sr.GetMyRequests(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	store := &calhandlers.StoredRequestsHandler{}
	http.HandleFunc("/api/stored", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Query().Get("slot_id") != "" {
			store.GetStoredRequestsForSlot(w, r)
		} else {
			store.GetStoredRequestsForTA(w, r)
		}
	})

	// Profile handlers — note the profhandlers alias
	bio := &profhandlers.BioHandler{Templates: tmpl}

	http.HandleFunc("/profiles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			bio.ProfilesPage(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ta", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			bio.TAPage(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/bios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bio.GetAllBiosJSON(w, r)
		case http.MethodPost:
			bio.UpsertBio(w, r)
		case http.MethodDelete:
			bio.DeleteBio(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	}

}

/*
# Create office hours slot
Invoke-RestMethod -Method POST http://localhost:8080/api/hours `
  -ContentType "application/json" `
  -Body '{"ta_name":"Andre","day":"2026-04-16","start_time":"17:30","end_time":"17:31","location":"Rice Hall 101"}'

# List your hours
Invoke-RestMethod -Method GET "http://localhost:8080/api/hours?ta_uid=uid-placeholder"

# Update — paste the id from the create response
Invoke-RestMethod -Method PATCH "http://localhost:8080/api/hours?id=YOUR_ID_HERE" `
  -ContentType "application/json" `
  -Body '{"ta_name":"Andre","day":"2026-04-16","start_time":"18:00","end_time":"19:00","location":"Rice Hall 202"}'

# Delete
Invoke-RestMethod -Method DELETE "http://localhost:8080/api/hours?id=YOUR_ID_HERE"

# Submit a student request — use the "id" from the create response as slot_id
Invoke-RestMethod -Method POST http://localhost:8080/api/requests `
  -ContentType "application/json" `
  -Body '{"slot_id":"PASTE_ID_HERE","ta_name":"Andre","day":"2026-04-16","start_time":"17:30","end_time":"17:31","location":"Rice Hall 101","reason":"Need help with HW3"}'

# TA views incoming requests
Invoke-RestMethod -Method GET "http://localhost:8080/api/requests?ta_uid=uid-placeholder"

# View requests for a specific slot
Invoke-RestMethod -Method GET "http://localhost:8080/api/requests?slot_id=YOUR_SLOT_ID_HERE"

# Student views their own requests
Invoke-RestMethod -Method GET "http://localhost:8080/api/requests/me?student_uid=student-uid-placeholder"

# Student cancels a request
Invoke-RestMethod -Method DELETE "http://localhost:8080/api/requests?id=YOUR_REQUEST_ID"

# TA views archived historical requests
Invoke-RestMethod -Method GET "http://localhost:8080/api/stored?ta_uid=uid-placeholder"
*/
