package handling

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"notif/internal/auth"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

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

type LoginRequest struct {
	Uid    string `json:"uid"`
	Exists bool   `json:"exists"`
}

var sessionCache = make(map[string]SessionCacheEntry) //session to id (has a longer expiration time)
var roleCache = make(map[string]RoleCacheEntry)       //id to roles (this is needed to ensure up-to-date roles to users)

// locking mechanisms to prevent read and write collision
var rc_lock sync.RWMutex
var sc_lock sync.RWMutex

// this needs to be called every time a request is made to validate the request
func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("reached middleware")
		fmt.Println("params: " + ps.ByName("page"))
		//session cookie comparison + creation
		cookie, err := r.Cookie(SESSION_COOKIE_NAME)
		if err != nil {
			//if cookie broken or not exist, make a new one!
			fmt.Print("cookie make")
			cookie, err = setSessionCookie(w)
			if err != nil {
				fmt.Println("1")
				loginRedirect(w, r)
				return
			}
		}
		if ps.ByName("page") == "login" { //pass through login if this is destination
			next(w, r, ps)
			return
		}

		//check for login needed (check sessioncache)
		sc_lock.RLock()
		var uid string
		if s, ex := sessionCache[cookie.Value]; !ex {
			sc_lock.RUnlock()
			fmt.Println("2")
			loginRedirect(w, r)
			return
		} else {
			uid = s.uid
		}

		sc_lock.RUnlock()

		//checking + updating cache
		rc_lock.RLock()
		if _, ex := roleCache[uid]; !ex {
			if roles, err := auth.GetRoles(uid); err != nil {
				log.Printf("Failed obtaining roles in authentication with error %v", err)
				rc_lock.RUnlock()
				fmt.Println("3")
				loginRedirect(w, r)
				return
			} else if roles == nil {
				rc_lock.RUnlock()
				fmt.Println("4")
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
		next(w, r, ps)
	}
}

func loginRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("reached login redirect")
	http.SetCookie(w, &http.Cookie{
		Name:     REDIRECT_COOKIE_NAME,
		Value:    url.QueryEscape(r.URL.RequestURI()),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})
	http.Redirect(w, r, "/view/login", http.StatusSeeOther)
}

// setting session cookie for connected client session
func setSessionCookie(w http.ResponseWriter) (*http.Cookie, error) {
	var cookie = http.Cookie{}
	sid, err := auth.GenerateSessionID()
	if err != nil {
		log.Printf("Failed session ID generation with error: %v", err)
		return &cookie, err
	}
	cookie = http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return &cookie, nil
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

func LoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("reached login handling")
	var v LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Printf("Failed decoding request json with error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uid := v.Uid
	exists := v.Exists
	fmt.Printf("uid: %v, exists: %v", uid, exists)
	if exists {
		roles, err := auth.GetRoles(uid)
		if err != nil {
			log.Printf("failed getting roles in login handling with error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	if cookie, err := r.Cookie(SESSION_COOKIE_NAME); err != nil {
		log.Printf("failed getting cookie in login handling with error: %v", err)
		http.Redirect(w, r, "/view/login", http.StatusFound)
		return
	} else {
		sessionCache[cookie.Value] = SessionCacheEntry{
			uid:        uid,
			expiration: time.Now().Add(60 * time.Minute),
		}
	}
	//redirect back to orginial page after login, or default to home
	var nurl = "/view/home"
	//only logging there being an actual issue, no cookie is to be expected if the user first reached the login through non-redirect means
	if cookie, err := r.Cookie(REDIRECT_COOKIE_NAME); err != nil && !errors.Is(err, http.ErrNoCookie) {
		log.Printf("Failed getting redirect cookie in login handling with error: %v", err)
	} else if cookie != nil {
		nurl, _ = url.QueryUnescape(cookie.Value)
	}
	fmt.Println("5")
	http.Redirect(w, r, nurl, http.StatusFound)
}

// simple function for comparing security and top role in user roles
func securityCheck(s int, rs []auth.RoleType) bool {
	var tr = rs[0]
	for _, r := range rs[1:] {
		if tr > r {
			tr = r
		}
	}
	fmt.Printf("\nroles: %v", rs)
	if s < int(tr) {
		fmt.Printf("\nFALSE; security: %v, toprole: %v", s, tr)
		return false
	} else {
		fmt.Printf("\nTRUE; security: %v, toprole: %v", s, tr)
		return true
	}
}

func grabUID(r *http.Request) (string, error) {
	if cookie, err := r.Cookie(SESSION_COOKIE_NAME); err != nil {
		return "", err
	} else {
		sc_lock.RLock()
		defer sc_lock.RUnlock()
		if entry, ex := sessionCache[cookie.Value]; !ex {
			return "", nil
		} else {
			return entry.uid, err
		}
	}
}
