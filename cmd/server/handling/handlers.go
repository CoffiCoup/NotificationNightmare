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
			//couldn't find roleentry
		} else {
			if !securityCheck(page.Security, e.roles) { //security check!
				//send the person to the gulag (front page, they shouldn't be here)
				//(LOOK INTO DIRECTING TO LAST PAGE OR THROWING UNAUTHORIZED WALL?)
			} else {
				http.ServeFile(w, r, page.URL)
			}
		}
	}
}
