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
	var v string
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("Failed to decode json from request body with error: %v", err)
		return
	}
	page := models.WEBPAGES[v]
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
	var v graph.OfficeHoursRow //updated office hour structure
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("Failed to decode json from request body with error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := storage.UpdateOfficeHoursJSON(v.ID, v); err != nil {
		log.Printf("Failed to update office hours json with error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent) //ok with nothing else
	}
}

func DeleteOHHandler(w http.ResponseWriter, r *http.Request) {
	var v string //id
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("Failed to decode json from request body with error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := storage.DeleteOfficeHoursJSON(v); err != nil {
		log.Printf("Failed to update office hours json with error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent) //ok with nothing else
	}
}

func CreateOHHandler(w http.ResponseWriter, r *http.Request) {
	var rOH graph.OfficeHoursRow //new office hours structure
	if err := json.NewDecoder(r.Body).Decode(&rOH); err != nil {
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
		w.WriteHeader(http.StatusCreated) //confirmed created
	}
}

func DeleteOHRHandler(w http.ResponseWriter, r *http.Request) {
	var v string //id
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("Failed to decode json from request body with error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := storage.DeleteStudentRequest(v); err != nil {
		log.Printf("Failed to update office hours json with error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent) //ok with nothing else
	}
}

func CreateOHRHandler(w http.ResponseWriter, r *http.Request) {
	var v graph.StudentRequest //id
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("Failed to decode json from request body with error: %v", err)
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
			v.StudentUID = entry.uid
		}
		sc_lock.RUnlock()
	}
	if _, err := storage.AppendStudentRequest(v); err != nil {
		log.Printf("Failed creating office hours through user input with error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated) //confirmed created
	}
}
