package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notif/internal/auth"
	"notif/internal/data"
	"notif/internal/models"
	"notif/internal/util/tests"
	"sync"
	"time"
)

// setting up constant key to label user_id obtained through reading http request (from netbadge)
// passed along through context to be used by handlers
type contextKey string

const USER_ID_KEY contextKey = "user_id"
const USER_ROLES contextKey = "user_role"
const SESSION_COOKIE_NAME string = "session_id"

// handling the cache
type RoleCacheEntry struct {
	roles      []auth.RoleType
	expiration time.Time
}

type SessionCacheEntry struct {
	uid        string
	expiration time.Time
}

var sessionCache = make(map[string]SessionCacheEntry) //session to id (has a longer expiration time)
var roleCache = make(map[string]RoleCacheEntry)       //id to roles (this is needed to ensure up-to-date roles to users)

// locking mechanisms to prevent read and write collision
var rc_lock sync.RWMutex
var sc_lock sync.RWMutex

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
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//session cookie comparison + creation
		cookie, err := r.Cookie(SESSION_COOKIE_NAME)
		if err != nil {
			sc_lock.Lock()
			//if cookie broken or not exist, make a new one!
			sid, err := auth.GenerateSessionID()
			if err != nil {
				log.Fatalf("failed to generate session id with error: %v", err)
			}
			http.SetCookie(w, &http.Cookie{
				Name:     SESSION_COOKIE_NAME,
				Value:    sid,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   3600,
			})
		}
		sc_lock.RLock()
		//login check FIX THIS!
		var uid string = sessionCache[cookie.Value].uid

		var ctx_v = r.Context().Value(USER_ID_KEY)
		if s, ok := ctx_v.(string); ok {
			uid = s
		} else {
			if ctx_v != nil {
				log.Fatalf("context obtained %v instead of string", ctx_v)
			}
		}
		if uid == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		sc_lock.Unlock()
		//FOR DEBUGGING, DELETE FOR TESTING CORRECT UID MESSAGING
		uid = "dev_user"
		if uid == "dev_user" {
			ctx := context.WithValue(r.Context(), USER_ID_KEY, uid)
			ctx = context.WithValue(ctx, USER_ROLES, []int{0})
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		//checking + updating cache
		rc_lock.RLock()
		var (
			entry RoleCacheEntry
			ex    bool
		)
		if entry, ex = roleCache[uid]; !ex {
			if roles := auth.GetRoles(uid); roles == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			} else {
				roleCache[uid] = RoleCacheEntry{
					roles:      roles,
					expiration: time.Now().Add(7 * time.Minute),
				}
			}
		}
		//adds user_id and roles to be read by handlers easily
		ctx := context.WithValue(r.Context(), USER_ID_KEY, uid)
		ctx = context.WithValue(ctx, USER_ROLES, entry.roles)
		next.ServeHTTP(w, r.WithContext(ctx))
		rc_lock.RUnlock()
	})
}

// running a goroutine to periodically clean out the cache
func cacheClean() {
	go func() {
		for {
			time.Sleep(7 * time.Minute)
			rc_lock.Lock()
			for uid, e := range roleCache {
				if time.Now().After(e.expiration) {
					delete(roleCache, uid)
				}
			}
			rc_lock.Unlock()
		}
	}()
	go func() {
		for {
			time.Sleep(15 * time.Minute)
			sc_lock.Lock()
			for sid, s := range sessionCache {
				if time.Now().After(s.expiration) {
					delete(sessionCache, sid)
				}
			}
			sc_lock.Unlock()
		}
	}()
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	uid := "temp dev" //Grab from login page request
	exists := true    //Grab from login page request
	if exists {
		roleCache[uid] = RoleCacheEntry{
			roles:      auth.GetRoles(uid),
			expiration: time.Now().Add(7 * time.Minute),
		}
	} else { //creating temp role for student
		//student doesn't need to be monitored NEARLY as much
		roleCache[uid] = RoleCacheEntry{
			roles:      []auth.RoleType{auth.Student},
			expiration: time.Now().Add(30 * time.Minute),
		}
	}
}

func main() {
	//server certificate and key file paths, make sure the files are actually here
	//these are too be used ONLY for debug, they will be removed upon sending this to production
	//these will need to be replaced with just a command line prompt when this is sent out...
	//servCert := "debugcerts/localhost.crt"
	//servKey := "debugcerts/localhost.key"

	//setting the default handlers for request signatures
	//TODO: view handler
	//TODO: oh update handler
	//TODO: o
	http.Handle("/admin", authMiddleware(http.HandlerFunc(tests.AdminPageHandler)))
	//handle login info from the login page (this needs to NOT be in the middleware!)
	http.HandleFunc("/login", loginHandler)

	cacheClean()
	//profiles.GetAllProfiles()

	//TESTING FUNCTION PLACEMENT HERE
	fmt.Println("tests starting")
	tests.RoleListTests()
	fmt.Println("\ntests ending")

	//tells the server to start listening to requests
	err := http.ListenAndServe("localhost:5500", nil)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
