package tests

import (
	"fmt"
	"net/http"
	"notif/internal/auth"
)

//TESTING AND HANDLIN TIME!

func AdminPageHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: Work on this with Ethan!
	fmt.Println("welcome admin")
}

func RoleListTests() {
	var entries = []auth.RoleEntry{
		{Uid: "001", Password: "password1", Roles: []auth.RoleType{1}},
		{Uid: "002", Password: "password2", Roles: []auth.RoleType{2}},
		{Uid: "003", Password: "password3", Roles: []auth.RoleType{3}},
	}
	var file = auth.MakeUpdateRoles(entries, nil)
	fmt.Printf("\nupdated roles made in %v", file.Name())
	auth.StoreRoleList()
	fmt.Println("\nstored current role list")
	auth.ReplaceRoleList(file)
	fmt.Println("replaced stored role list")
	fmt.Printf("\nUser 001, Roles: %v", auth.GetRoles("001"))
	fmt.Printf("\nUser 002, Roles: %v", auth.GetRoles("002"))
	fmt.Printf("\nUser 003, Roles: %v", auth.GetRoles("003"))
	fmt.Printf("\nUser 001, Password: %v", auth.GetPassword("001"))
	fmt.Printf("\nUser 002, Password: %v", auth.GetPassword("002"))
	fmt.Printf("\nUser 003, Password: %v", auth.GetPassword("003"))
}
