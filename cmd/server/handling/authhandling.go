package handling

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"notif/internal/auth"
	"notif/internal/models"
	"sync"
	"time"
)

// setting up constant key to label user_id obtained through reading http request (from netbadge)
// passed along through context to be used by handlers
type contextKey string

const USER_ID_KEY contextKey = "user_id"
const USER_ROLES contextKey = "user_role"
const SESSION_COOKIE_NAME string = "session_id"
const REDIRECT_COOKIE_NAME string = "redirect_url"

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

// this needs to be called every time a request is made to validate the request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//session cookie comparison + creation
		var sid string
		cookie, err := r.Cookie(SESSION_COOKIE_NAME)
		if err != nil {
			//if cookie broken or not exist, make a new one!
			sid, err = auth.GenerateSessionID()
			if err != nil {
				log.Printf("Failed to generate session id with error: %v", err)
				loginRedirect(w, r)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     SESSION_COOKIE_NAME,
				Value:    sid,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   3600,
			})
		} else {
			sid = cookie.Value
		}

		//check for login needed (check sessioncache)
		sc_lock.RLock()
		var uid string
		if s, ex := sessionCache[sid]; !ex {
			sc_lock.RUnlock()
			loginRedirect(w, r)
			return
		} else {
			uid = s.uid
		}

		sc_lock.RUnlock()

		// //FOR DEBUGGING, DELETE FOR TESTING CORRECT UID MESSAGING
		// uid = "dev_user"
		// if uid == "dev_user" {
		// 	ctx := context.WithValue(r.Context(), USER_ID_KEY, uid)
		// 	ctx = context.WithValue(ctx, USER_ROLES, []int{0})
		// 	next.ServeHTTP(w, r.WithContext(ctx))
		// 	return
		// }

		//checking + updating cache
		rc_lock.RLock()
		var (
			entry RoleCacheEntry
			ex    bool
		)
		if entry, ex = roleCache[uid]; !ex {
			if roles, err := auth.GetRoles(uid); err != nil {
				log.Printf("Failed obtaining roles in authentication with error %v", err)
				rc_lock.RUnlock()
				loginRedirect(w, r)
				return
			} else if roles == nil {
				rc_lock.RUnlock()
				loginRedirect(w, r)
				return
			} else {
				roleCache[uid] = RoleCacheEntry{
					roles:      roles,
					expiration: time.Now().Add(7 * time.Minute),
				}
			}
		}
		rc_lock.RUnlock()
		//adds user_id and roles to be read by handlers easily
		ctx := context.WithValue(r.Context(), USER_ID_KEY, uid)
		ctx = context.WithValue(ctx, USER_ROLES, entry.roles)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loginRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, models.WEBPAGES["login"].URL, http.StatusSeeOther)
	http.SetCookie(w, &http.Cookie{
		Name:     REDIRECT_COOKIE_NAME,
		Value:    url.QueryEscape(r.URL.RequestURI()),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})

}

// running a goroutine to periodically clean out the cache
func CacheClean() {
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: have ethan help me here
	uid := "temp dev" //Grab from login page request
	exists := true    //Grab from login page request
	if exists {
		roles, err := auth.GetRoles(uid)
		if err != nil {
			log.Printf("failed getting roles in login handling with error: %v", err)
			loginRedirect(w, r)
			return
		}
		roleCache[uid] = RoleCacheEntry{
			roles:      roles,
			expiration: time.Now().Add(7 * time.Minute),
		}
	} else { //creating temp role for student
		//student role doesn't need to be monitored NEARLY as much
		roleCache[uid] = RoleCacheEntry{
			roles:      []auth.RoleType{auth.Student},
			expiration: time.Now().Add(60 * time.Minute),
		}
	}
	//storing session to uid in cache
	if cookie, err := r.Cookie(SESSION_COOKIE_NAME); err != nil { //if the cookie can't be found, something went wrong so just direct to home
		log.Printf("failed getting cookie in login handling with error: %v", err)
		http.Redirect(w, r, models.WEBPAGES["home"].URL, http.StatusFound)
		return
	} else {
		sessionCache[cookie.Value] = SessionCacheEntry{
			uid:        uid,
			expiration: time.Now().Add(60 * time.Minute),
		}
	}
	//redirect back to orginial page after login, or default to home
	var nurl = models.WEBPAGES["home"].URL
	//only logging there being an actual issue, no cookie is to be expected if the user first reached the login through non-redirect means
	if cookie, err := r.Cookie(REDIRECT_COOKIE_NAME); err != nil && !errors.Is(err, http.ErrNoCookie) {
		log.Printf("Failed getting redirect cookie in login handling with error: %v", err)
	} else if cookie != nil {
		nurl, _ = url.QueryUnescape(cookie.Value)
	}
	http.Redirect(w, r, nurl, http.StatusFound)
}

// simple function for comparing security and top role in user roles
func securityCheck(s int, rs []auth.RoleType) bool {
	var tr = rs[0]
	for _, r := range rs[1:] {
		if tr < r {
			tr = r
		}
	}
	if s < int(tr) {
		return false
	} else {
		return true
	}
}
