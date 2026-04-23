package handling

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"notif/internal/auth"
	"notif/internal/calendar/graph"
	"notif/internal/profiles/storage"

	"github.com/julienschmidt/httprouter"
)

type BioHandler struct {
	Templates *template.Template
}

func GETProfileCentralHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var biohandler = BioHandler{}
	switch ps.ByName("action") {
	case "grid":
		biohandler.ProfilesGridHandler(w, r, ps)
	case "editor":
		biohandler.BioEditorHandler(w, r)
	}
}

func POSTProfileCentralHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var biohandler = BioHandler{}
	switch ps.ByName("action") {
	case "upsert":
		biohandler.upsertBioHandler(w, r)
	case "delete":
		biohandler.deleteBioHandler(w, r)
	}
}

// ProfilesPage handles GET /profiles
// Renders the public read-only grid of all TA bios
func (h *BioHandler) ProfilesGridHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bios, err := storage.GetAllBios()
	if err != nil {
		log.Println("ERROR fetching bios:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	if err := h.Templates.ExecuteTemplate(w, "profiles.html", bios); err != nil {
		log.Println("ERROR rendering profiles template:", err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// TAPage handles GET /ta
// Renders the TA dashboard where they can create/edit their bio
func (h *BioHandler) BioEditorHandler(w http.ResponseWriter, r *http.Request) {
	var taUID string
	if uid, err := grabUID(r); err != nil { //grabbing uid
		log.Printf("Failed grabbing uid with error: %v", err)
		loginRedirect(w, r)
		return
	} else if uid == "" {
		loginRedirect(w, r)
		return
	} else {
		taUID = uid
	}

	if roles, err := auth.GetRoles(taUID); err != nil { //security
		log.Printf("Failed obtaining roles with error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if !securityCheck(2, roles) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	// Try to load existing bio to pre-fill the form
	bio, err := storage.GetBioByTA(taUID)
	if err != nil {
		// No bio yet — pass empty struct so template still renders
		bio = &graph.TABio{}
	}

	if err := h.Templates.ExecuteTemplate(w, "ProfileUpload.html", bio); err != nil {
		log.Println("ERROR rendering TA template:", err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

// UpsertBio handles POST /api/bios
// Called when TA submits the form on /ta
func (h *BioHandler) upsertBioHandler(w http.ResponseWriter, r *http.Request) {
	var bio graph.TABio
	if err := json.NewDecoder(r.Body).Decode(&bio); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if bio.Name == "" || bio.Email == "" {
		http.Error(w, "name and email are required", http.StatusBadRequest)
		return
	}

	if uid, err := grabUID(r); err != nil { //grabbing uid
		log.Printf("Failed grabbing uid with error: %v", err)
		loginRedirect(w, r)
		return
	} else if uid == "" {
		loginRedirect(w, r)
		return
	} else {
		bio.TAUID = uid
	}

	id, err := storage.UpsertBio(bio)
	if err != nil {
		log.Println("ERROR saving bio:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     id,
	})
}

// DeleteBio handles DELETE /api/bios
func (h *BioHandler) deleteBioHandler(w http.ResponseWriter, r *http.Request) {
	var taUID string
	if uid, err := grabUID(r); err != nil { //grabbing uid
		log.Printf("Failed grabbing uid with error: %v", err)
		loginRedirect(w, r)
		return
	} else if uid == "" {
		loginRedirect(w, r)
		return
	} else {
		taUID = uid
	}

	if taUID == "" {
		taUID = "BLANK"
	}

	if err := storage.DeleteBio(taUID); err != nil {
		log.Println("ERROR deleting bio:", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllBiosJSON handles GET /api/bios
// Returns JSON — useful for future JavaScript calls
func (h *BioHandler) GetAllBiosJSON(w http.ResponseWriter, r *http.Request) {
	bios, err := storage.GetAllBios()
	if err != nil {
		log.Println("ERROR fetching bios:", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bios)
}

// Renders a single TA's full profile page
func (h *BioHandler) ProfilePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	taUID := ps.ByName("uid")
	if taUID == "" {
		http.Error(w, "missing uid", http.StatusBadRequest)
		return
	}

	bio, err := storage.GetBioByTA(taUID)
	if err != nil {
		log.Println("ERROR fetching bio:", err)
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	if err := h.Templates.ExecuteTemplate(w, "profile.html", bio); err != nil {
		log.Println("ERROR rendering profile template:", err)
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}
