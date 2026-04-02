package auth

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"slices"
	"sync"
)

//security and authentication related functions should be handled here (outside of the basics handled in main.go)

const ROLELISTPATH string = "internal/auth/roleList.json"

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
	uid   string
	roles []string
}

func (r RoleType) String() string {
	return roleTypeName[r]
}

// handling the cache
var roleCache = make(map[string]RoleEntry)
var mu sync.RWMutex

func getRole(uid string) []string {
	mu.RLock()
	roleListFile, err := os.Open(ROLELISTPATH)
	if err != nil {
		log.Fatal("Failed roleList.json Read")
		mu.RUnlock()
		return nil
	}
	defer roleListFile.Close()
	decoder := json.NewDecoder(roleListFile)
	for decoder.More() {
		var v RoleEntry
		if err := decoder.Decode(&v); err != nil {
			log.Println("roleList decode error")
			mu.RUnlock()
			return nil
		}
		if v.uid == uid {
			defer mu.RUnlock()
			return v.roles
		}
	}
	log.Fatal("failed to get role")
	return nil
}

func makeUpdateRoles(entries []RoleEntry) *os.File {
	mu.RLock()
	defer mu.RUnlock()
	lookupEntries := map[string]RoleEntry{}
	for _, e := range entries {
		lookupEntries[e.uid] = e
	}
	roleListNewFile, err := os.Create(path.Join(ROLELISTPATH, "_"))
	if err != nil {
		log.Fatalf("Failed %v creation", path.Join(ROLELISTPATH, "_"))
	}
	roleListFile, err := os.Open(ROLELISTPATH)
	if err != nil {
		log.Fatal("Failed roleList.json Read")
		return nil
	}
	defer roleListFile.Close()
	decoder := json.NewDecoder(roleListFile)
	encoder := json.NewEncoder(roleListNewFile)
	//update loop
	for decoder.More() {
		var v RoleEntry
		if err := decoder.Decode(&v); err != nil {
			log.Println("roleList decode error")
			return roleListNewFile
		}
		if e, ok := lookupEntries[v.uid]; ok {
			if r := e.roles; !slices.Equal(r, v.roles) {
				v.roles = r
			}
		} else {
			//add to a list to append to everything at the end
		}
		if err := encoder.Encode(&v); err != nil {
			log.Println("roleList encode error")
			return roleListNewFile
		}
	}
	//create loop
	return roleListNewFile
}

func updateRolelist(newRoleList *os.File) {
	newRoleList.Name()

}
