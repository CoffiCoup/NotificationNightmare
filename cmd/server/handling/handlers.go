package handling

import (
	"encoding/json"
	"log"
	"net/http"
	"notif/internal/calendar/graph"
	"notif/internal/calendar/storage"
	"notif/internal/models"
)

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: get ethan to help me with the request here
	pn := "login" //page name
	page := models.WEBPAGES[pn]
	if cookie, err := r.Cookie(SESSION_COOKIE_NAME); err != nil {
		log.Printf("Failed to obtain session cookie with error: %v", err)
	} else {
		if e, ex := roleCache[cookie.Value]; !ex {
			loginRedirect(w, r)
			return
		} else {
			if !securityCheck(page.Security, e.roles) { //security check!
				http.Error(w, "Unauthorized", http.StatusForbidden)
				return
			} else { //if all goes well, serve it
				http.ServeFile(w, r, page.URL)
			}
		}
	}
}

func UpdateOHHandler(w http.ResponseWriter, r *http.Request) {
	id := "" //obtained from Ethan communication
	day := ""
	start := ""
	end := ""
	location := ""
	storage.UpdateOfficeHoursJSON(id, day, start, end, location)
}

func DeleteOHHandler(w http.ResponseWriter, r *http.Request) {
	id := "" //obtained from Ethan communication
	storage.DeleteOfficeHoursJSON(id)
}

func CreateOHHandler(w http.ResponseWriter, r *http.Request) {
	rBody := r.Body
	var rOH graph.OfficeHoursRow
	//get ta Name from ethan request html
	if err := json.NewDecoder(rBody).Decode(&rOH); err != nil {
		log.Printf("Failed to decode request body json with error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if cookie, err := r.Cookie(SESSION_COOKIE_NAME); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		loginRedirect(w, r)
		return
	} else {
		sc_lock.RLock()
		if entry, ex := sessionCache[cookie.Value]; !ex {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			loginRedirect(w, r)
			sc_lock.RUnlock()
			return
		} else {
			rOH.TAUID = entry.uid
		}
		sc_lock.RUnlock()
	}
	if _, err := storage.AppendOfficeHoursJSON(rOH); err != nil {
		log.Printf("Failed creating office hours through user input with error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteOHRHandler(w http.ResponseWriter, r *http.Request) {
	//do
}

func CreateOHRHandler(w http.ResponseWriter, r *http.Request) {
	//do
}
