package main

import (
	"fmt"
	"log"
	"net/http"
	"notif/cmd/server/handling"
	"notif/internal/util/tests"
)

func main() {
	//server certificate and key file paths, make sure the files are actually here
	//these are too be used ONLY for debug, they will be removed upon sending this to production
	//these will need to be replaced with just a command line prompt when this is sent out...
	//servCert := "debugcerts/localhost.crt"
	//servKey := "debugcerts/localhost.key"

	//setting the default handlers for request signatures
	http.Handle("/view", handling.AuthMiddleware(http.HandlerFunc(handling.ViewHandler)))
	http.HandleFunc("/ohcreate", handling.CreateOHHandler)
	http.HandleFunc("/ohupdate", handling.UpdateOHHandler)
	http.HandleFunc("/ohdelete", handling.DeleteOHHandler)
	http.HandleFunc("/ohrcreate", handling.CreateOHRHandler)
	http.HandleFunc("/ohrdelete", handling.DeleteOHRHandler)
	//handle login info from the login page (this needs to NOT be in the middleware!)
	http.HandleFunc("/login", handling.LoginHandler)

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
