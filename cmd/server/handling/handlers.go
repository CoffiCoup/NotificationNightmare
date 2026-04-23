package handling

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notif/internal/calendar/graph"
	"notif/internal/calendar/storage"
	"notif/internal/models"
	"os"

	"github.com/julienschmidt/httprouter"
)

func ViewHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pn := ps.ByName("page")
	var page, ex = models.WEBPAGES[pn]
	if !ex {
		http.NotFound(w, r)
		return
	} else if pn == "login" { //skip to login if that's destination
		http.ServeFile(w, r, page.URL)
		return
	}
	if cookie, err := r.Cookie(SESSION_COOKIE_NAME); err != nil {
		log.Printf("Failed to obtain session cookie with error: %v", err)
	} else {
		if es, ex := sessionCache[cookie.Value]; !ex {
			fmt.Println("7")
			loginRedirect(w, r, ps)
			return
		} else {
			if er, ex := roleCache[es.uid]; !ex {
				fmt.Println("6")
				loginRedirect(w, r, ps)
				return
			} else {
				if !securityCheck(page.Security, er.roles) { //security check!
					http.Error(w, "Unauthorized", http.StatusForbidden)
					return
				} else { //if all goes well, serve it
					http.ServeFile(w, r, page.URL)
				}
			}
		}
	}
}

func OHCentralHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	switch action {
	case "create":
		CreateOHHandler(w, r, ps)
	case "update":
		UpdateOHHandler(w, r)
	case "delete":
		DeleteOHHandler(w, r, ps)
	default:
		http.NotFound(w, r)
	}
}

func OHRCentralHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	switch action {
	case "create":
		CreateOHRHandler(w, r, ps)
	case "delete":
		DeleteOHRHandler(w, r, ps)
	default:
		http.NotFound(w, r)
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

func DeleteOHHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id string = ps.ByName("extra")
	if err := storage.DeleteOfficeHoursJSON(id); err != nil {
		log.Printf("Failed to update office hours json with error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent) //ok with nothing else
	}
}

func CreateOHHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var rOH graph.OfficeHoursRow //new office hours structure
	if err := json.NewDecoder(r.Body).Decode(&rOH); err != nil {
		log.Printf("Failed to decode request body json with error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if cookie, err := r.Cookie(SESSION_COOKIE_NAME); err != nil {
		loginRedirect(w, r, ps)
		return
	} else {
		sc_lock.RLock()
		if entry, ex := sessionCache[cookie.Value]; !ex {
			sc_lock.RUnlock()
			loginRedirect(w, r, ps)
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

func DeleteOHRHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var id string = ps.ByName("extra")
	if err := storage.DeleteStudentRequest(id); err != nil {
		log.Printf("Failed to update office hours json with error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent) //ok with nothing else
	}
}

func CreateOHRHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var v graph.StudentRequest //something something
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("Failed to decode json from request body with error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if cookie, err := r.Cookie(SESSION_COOKIE_NAME); err != nil {
		loginRedirect(w, r, ps)
		return
	} else {
		sc_lock.RLock()
		if entry, ex := sessionCache[cookie.Value]; !ex {
			sc_lock.RUnlock()
			loginRedirect(w, r, ps)
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

func FetchCentralHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	switch ps.ByName("file") {
	case "rolelist":
		roleListFetchHandler(w, r, ps)
	default:
		http.NotFound(w, r)
	}

}

func roleListFetchHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	file, err := os.ReadFile("internal/auth/roleList.json")
	if err != nil {
		log.Printf("Failed to obtain data from rolelist with error: %v", err)
		http.Error(w, "file not found", http.StatusInternalServerError)
		return
	}
	w.Write(file)
}
