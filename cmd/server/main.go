package main

import (
	"context"
	"log"
	"net/http"
	"notif/internal/data"   // Your data package
	"notif/internal/models" // Your models package
	"notif/cmd/server/handlers"
)

// setting up constant key to label user_id obtained through reading http request (from netbadge)
// passed along through context to be used by handlers
type contextKey string

const USER_ID_KEY contextKey = "user_id"

//SAMANVI01 CHANGES TO MAIN

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

//END SAMANVI01 CHANGES TO MAIN

// processing the login data from netbadge log-in (log-in handled by the hosting server module)
// this needs to be called every time a request is made to validate the request
func netBadgeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := r.Header.Get("HTTP_UID")

		//FOR DEBUGGING, DELETE FOR TESTING CORRECT UID MESSAGING
		uid = "dev_user"

		//evaluating uid, will need to implement roles here later to attach to permission levels
		if uid == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		//TODO: add to cookie user role

		//adds user_id key to context to pass around
		ctx := context.WithValue(r.Context(), USER_ID_KEY, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	//server certificate and key file paths, make sure the files are actually here
	//these are too be used ONLY for debug, they will be removed upon sending this to production
	//these will need to be replaced with just a command line prompt when this is sent out...
	//servCert := "debugcerts/localhost.crt"
	//servKey := "debugcerts/localhost.key"

	//setting the default handlers for request signatures
	http.HandleFunc("/view", handlers.StudentViewHandler)
  // New Excel Integration Endpoint
	http.HandleFunc("/signup", signupHandler)
	//tells the server to start listening to requests
	err := http.ListenAndServe("5500", nil)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
