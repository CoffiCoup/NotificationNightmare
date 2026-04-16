package main

import (
	"fmt"
	"log"
	"net/http"
	"notif/cmd/server/handling"
	"notif/internal/data"
	"notif/internal/models"
	"notif/internal/util/tests"
)

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

func main() {
	//server certificate and key file paths, make sure the files are actually here
	//these are too be used ONLY for debug, they will be removed upon sending this to production
	//these will need to be replaced with just a command line prompt when this is sent out...
	//servCert := "debugcerts/localhost.crt"
	//servKey := "debugcerts/localhost.key"

	//setting the default handlers for request signatures
	//TODO: view handler
	//TODO: oh update handler
	//TODO: oh request handler
	http.Handle("/admin", handling.AuthMiddleware(http.HandlerFunc(tests.AdminPageHandler)))
	//handle login info from the login page (this needs to NOT be in the middleware!)
	http.HandleFunc(models.WEBPAGES["login"].URL, handling.LoginHandler)

	handling.CacheClean()
	//profiles.GetAllProfiles()

	//TESTING FUNCTION PLACEMENT HERE
	fmt.Println("tests starting")
	// tests.RoleListTests()
	tests.TestCentral()
	fmt.Println("\ntests ending")

	//tells the server to start listening to requests
	err := http.ListenAndServe("localhost:5500", nil)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
