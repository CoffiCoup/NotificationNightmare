package handling

import (
	"log"
	"net/http"
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
