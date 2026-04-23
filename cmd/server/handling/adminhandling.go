package handling

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notif/internal/auth"
	"os"

	"github.com/julienschmidt/httprouter"
)

func GETAdminCentralHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	switch action {
	case "getpassword":
		passRolePasswordHandler(w, ps)
	default:
		http.NotFound(w, r)
	}
}

func POSTAdminCentralHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	switch action {
	case "updateroles":
		updateRolesHandler(w, r)
	case "restoreroles":
		restoreRolesHandler(w)
	default:
		http.NotFound(w, r)
	}
}

func updateRolesHandler(w http.ResponseWriter, r *http.Request) {
	var v [][]auth.RoleEntry
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("Failed decoding json from request for update roles with error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var file *os.File
	if f, err := auth.MakeUpdateRoles(v[0], v[1]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		file = f
		return
	}
	if err := auth.StoreRoleList(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := auth.ReplaceRoleList(file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func restoreRolesHandler(w http.ResponseWriter) {
	if err := auth.RestoreRoleList(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func passRolePasswordHandler(w http.ResponseWriter, ps httprouter.Params) {
	uid := ps.ByName("extra") //uid from request
	if pawd, err := auth.GetPassword(uid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Printf("\n uid: %v, pawd: %v", uid, pawd)
		w.Write([]byte(pawd))
	}
}
