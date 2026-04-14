package main

import (
	"log"
	"net/http"

	"notif/internal/calendar/graph"
	"notif/internal/calendar/handlers"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	token, err := graph.GetToken()
	if err != nil {
		log.Fatal("FATAL: could not get app token:", err)
	}
	log.Println("App token acquired:", token[:20]+"...")

	// Office hours handlers (TA)
	h := &handlers.HoursHandler{}

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

	// Student request handlers
	sr := &handlers.StudentReqHandler{}

	http.HandleFunc("/api/requests", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			sr.SubmitRequest(w, r)
		case http.MethodGet:
			// Route based on query param
			if r.URL.Query().Get("event_id") != "" {
				sr.GetRequestsForEvent(w, r)
			} else {
				sr.GetRequestsForTA(w, r)
			}
		case http.MethodDelete:
			sr.DeleteRequest(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Separate route for students viewing their own requests
	http.HandleFunc("/api/requests/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			sr.GetMyRequests(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}

/*
//Test creation, deletion, modification, and listing of office hours
# Create — note the "id" in the response
Invoke-RestMethod -Method POST http://localhost:8080/api/hours `
  -ContentType "application/json" `
  -Body '{"ta_name":"Andre","day":"2026-04-15","start_time":"14:00","end_time":"15:00","location":"Rice Hall 101"}'

# List your hours
Invoke-RestMethod -Method GET "http://localhost:8080/api/hours?ta_uid=uid-placeholder"

# Update — paste the id from the create response
Invoke-RestMethod -Method PATCH "http://localhost:8080/api/hours?id=YOUR_ID_HERE" `
  -ContentType "application/json" `
  -Body '{"ta_name":"Andre","day":"2026-04-14","start_time":"15:00","end_time":"16:00","location":"Rice Hall 202"}'

# Delete
Invoke-RestMethod -Method DELETE "http://localhost:8080/api/hours?id=YOUR_ID_HERE"
*/

/*
Test viewing and canceling student requests
# First create an office hours slot and note the outlook_event_id from officehours.json
# Then submit a student request using that event ID:
Invoke-RestMethod -Method POST http://localhost:8080/api/requests `
  -ContentType "application/json" `
  -Body '{
    "outlook_event_id": "YOUR_OUTLOOK_EVENT_ID_HERE",
    "ta_name": "Andre",
    "day": "2026-04-15",
    "start_time": "14:00",
    "end_time": "15:00",
    "location": "Rice Hall 101",
    "reason": "Need help with homework 3 question 4"
  }'

# TA views all incoming requests for their slots
Invoke-RestMethod -Method GET "http://localhost:8080/api/requests?ta_uid=uid-placeholder"

# View requests for a specific event
Invoke-RestMethod -Method GET "http://localhost:8080/api/requests?event_id=YOUR_OUTLOOK_EVENT_ID_HERE"

# Student views their own requests
Invoke-RestMethod -Method GET "http://localhost:8080/api/requests/me?student_uid=student-uid-placeholder"

# Student cancels a request
Invoke-RestMethod -Method DELETE "http://localhost:8080/api/requests?id=YOUR_REQUEST_ID"
*/
