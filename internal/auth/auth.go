package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"slices"
	"sync"
)

//security and authentication related functions should be handled here (outside of the basics handled in main.go)

const ROLELISTPATH string = "internal/auth/roleList.json"
const ROLELISTDIR string = "internal/auth/"

// role values and structures for authentication
type RoleType int

const (
	Admin RoleType = iota
	Staff
	TA
	Student
	Guest
)

var roleTypeName = map[RoleType]string{
	Admin:   "Admin",
	Staff:   "Staff",
	TA:      "TA",
	Student: "Student", //default, unless netbadge can implement all of this properly...
	Guest:   "Guest",
}

type RoleEntry struct {
	Uid      string
	Password string
	Roles    []RoleType
}

func (r RoleType) String() string {
	return roleTypeName[r]
}

var mu sync.RWMutex

func GetPassword(uid string) string {
	//file setup
	mu.RLock()
	defer mu.RUnlock()
	roleListFile, err := os.Open(ROLELISTPATH)
	if err != nil {
		log.Print("Failed roleList.json read for password")
		return ""
	}
	defer roleListFile.Close()
	//decoding rolelist
	decoder := json.NewDecoder(roleListFile)
	var e []RoleEntry
	if err := decoder.Decode(&e); err != nil {
		log.Printf("password roleList decode error: %v", err)
		return ""
	}
	//searching for password
	for _, v := range e {
		if v.Uid == uid {
			return v.Password
		}
	}
	return ""
}

func GetRoles(uid string) []RoleType {
	//file setup
	mu.RLock()
	defer mu.RUnlock()
	roleListFile, err := os.Open(ROLELISTPATH)
	if err != nil {
		log.Print("Failed roleList.json read for roles")
		return nil
	}
	defer roleListFile.Close()
	//decoding rolelist
	decoder := json.NewDecoder(roleListFile)
	var e []RoleEntry
	if err := decoder.Decode(&e); err != nil {
		log.Printf("roles roleList decode error: %v", err)
		return nil
	}
	//searching for password
	for _, v := range e {
		if v.Uid == uid {
			return v.Roles
		}
	}
	return nil
}

// creates an updated rolelist file, entries is updates or new entries, rEntries is entries to be removed
func MakeUpdateRoles(entries []RoleEntry, rentries []RoleEntry) *os.File {
	fmt.Print(entries)
	mu.RLock()
	defer mu.RUnlock()
	//creating lookup maps for comparison
	lookupEntries := map[string]RoleEntry{}
	rlookupEntries := map[string]RoleEntry{}
	for _, e := range entries {
		lookupEntries[e.Uid] = e
	}
	for _, e := range rentries {
		rlookupEntries[e.Uid] = e
	}
	//file prep
	roleListNewFile, err := os.CreateTemp(ROLELISTDIR, "roleList_new_*.json")
	if err != nil {
		log.Printf("Failed %v creation", path.Join(ROLELISTPATH, "_"))
	}
	roleListFile, err := os.Open(ROLELISTPATH)
	var ogFileExists bool = true
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			ogFileExists = false
		} else {
			log.Printf("Failed rolelist.json read with %v", err)
			return nil
		}
	}
	defer roleListFile.Close()
	decoder := json.NewDecoder(roleListFile)
	encoder := json.NewEncoder(roleListNewFile)
	defer roleListNewFile.Close()
	//grabbing current roles
	finalEntries := map[string]RoleEntry{}
	if ogFileExists {
		var e []RoleEntry
		if err := decoder.Decode(&e); err != nil {
			log.Printf("rolelist roleList decode error: %v", err)
			return roleListNewFile
		}
		for _, v := range e {
			if _, ex := rlookupEntries[v.Uid]; ex { //removal by skipping adding to final
				continue
			}
			finalEntries[v.Uid] = v
		}
	}
	//adding on the changes to the final entries
	for uid, v := range lookupEntries {
		if e, ex := finalEntries[uid]; ex { //update entries
			if r := e.Roles; !slices.Equal(r, v.Roles) {
				e.Roles = v.Roles
			}
			if p := e.Password; !(p == v.Password) {
				e.Password = v.Password
			}
		} else { //new entries
			finalEntries[uid] = v
		}
	}
	//encode finalentries to new file (formats first to make json work)
	var finalFmtEntries = make([]RoleEntry, 0, len(finalEntries))
	for _, e := range finalEntries {
		finalFmtEntries = append(finalFmtEntries, e)
	}
	if err := encoder.Encode(finalFmtEntries); err != nil {
		log.Println("failed to encode updated role list to new json file")
	}
	return roleListNewFile
}

// takes the temp rolelist supplied and makes it the main rolelist;
// should run storeRoleList() first in most cases because this does delete the og file
func ReplaceRoleList(newRoleList *os.File) {
	if err := os.Remove(ROLELISTPATH); err != nil { //removing current file
		if !errors.Is(err, fs.ErrNotExist) { //error only if error is not "file doesn't exist"
			log.Println("failed to remove current role list file")
		}
	}
	if err := os.Rename(newRoleList.Name(), ROLELISTPATH); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Println("attempted to replace role list with nonexistent file ")
		} else {
			log.Printf("\nfailed to replace old list file with new list file with error %v", err)
		}
	}
}

// backing up the second most recent list in case any errors occured
func StoreRoleList() {
	if err := os.Remove(ROLELISTDIR + "roleList_store.json"); err != nil { //removing current file
		if !errors.Is(err, fs.ErrNotExist) { //error only if error is not "file doesn't exist"
			log.Println("failed to remove current stored rolelist file")
		}
	}
	if err := os.Rename(ROLELISTPATH, ROLELISTDIR+"roleList_store.json"); err != nil { //changing filepath to storage path
		if errors.Is(err, fs.ErrNotExist) {
			log.Println("attempted to store nonexistent rolelist")
		} else {
			log.Printf("\nfailed to store roleList file with error %v", err)
		}
	}
}

// restoring backup
func RestoreRoleList() {
	if err := os.Remove(ROLELISTPATH); err != nil { //removing current file
		if !errors.Is(err, fs.ErrNotExist) { //error only if error is not "file doesn't exist"
			log.Println("failed to remove current stored rolelist file")
		}
	}
	if err := os.Rename((ROLELISTDIR + "roleList_store.json"), ROLELISTPATH); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Println("attempted to restore nonexistent stored role list")
		} else {
			log.Printf("\nfailed to restore stored roleList with error %v", err)
		}
	}
}

// generating encrypted sessionid
func GenerateSessionID() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// login function
func authorizeUser(uid string, password string)
