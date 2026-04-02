package main

import (
	"fmt"
	"log"
	"net/http"
	"notif/internal/data"   // Your data package
	"notif/internal/models" // Your models package
)

// signupHandler: Aligned with your existing models.OHRequest
func signupHandler(w http.ResponseWriter, r *http.Request) {
	// Using the exact field names from your models: ComputingID, TAID, DateTime, Reason
	newReq := models.OHRequest{
		ComputingID: "mst3k", // The Student
		TAID:        "samanvi_01",
		DateTime:    "2026-03-27 14:00",
		Reason:      "Help with Excel integration",
	}

	// Calling your existing SaveOHRequest function
	// We pass the filename we want to save to
	err := data.SaveOHRequest("OH_Requests.xlsx", newReq)
	if err != nil {
		log.Printf("Error saving to Excel: %v", err)
		http.Error(w, "Failed to save request", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully booked OH for %s with TA %s!", newReq.ComputingID, newReq.TAID)
}

func examplehandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the OH System at %s!", r.URL.Path[1:])
}

func certificationRequired(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Client certificate required", http.StatusUnauthorized)
		}
	}
}

func main() {
	// Standard Endpoints
	http.HandleFunc("/", examplehandler)

	// New Excel Integration Endpoint
	http.HandleFunc("/signup", signupHandler)

	// Secure Endpoint
	http.HandleFunc("/secure", certificationRequired(http.HandlerFunc(examplehandler)))

	certFile := "cert.crt"
	keyFile := "private.key"

	fmt.Println("Server starting on https://localhost:8080...")

	// Reminder: Ensure cert.crt and private.key are in your project folder
	err := http.ListenAndServeTLS(":8080", certFile, keyFile, nil)
	if err != nil {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}
