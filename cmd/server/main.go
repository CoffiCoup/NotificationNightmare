package main

//the starting point for the project for everything

import (
	"context"
	"log"
	"net/http"
	"notif/cmd/server/handlers"
)

// setting up constant key to label user_id obtained through reading http request (from netbadge)
// passed along through context to be used by handlers
type contextKey string

const USER_ID_KEY contextKey = "user_id"

// REMINDER: used ai to help find this function, add to notebook
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
	//tells the server to start listening to requests
	err := http.ListenAndServe("5500", nil)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
