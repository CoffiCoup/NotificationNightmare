package main

import (
	"html/template"
	"log"
	"net/http"
	"notif/cmd/server/handling"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Load templates from profiles package
	tmpl, err := template.ParseFiles(
		"internal/pages/profile.html",
		"internal/pages/profiles.html",
		"internal/pages/ProfileUpload.html",
	)
	if err != nil {
		log.Fatal("FATAL: could not load templates:", err)
	}

	var h = handling.BioHandler{Templates: tmpl}

	router := httprouter.New()
	router.GET("/view/:page", handling.AuthMiddleware(handling.ViewHandler)) //viewing
	router.GET("/fetch/:file", handling.FetchCentralHandler)                 //sending files to html
	router.GET("/profile/:action/:extra", h.GETProfileCentralHandler)        //GET profile requests
	router.GET("/admin/:action/:extra", handling.GETAdminCentralHandler)     //GET admin requests (rolelist, etc.)
	router.POST("/oh/:action/:extra", handling.OHCentralHandler)             //office hour stuff
	router.POST("/ohr/:action", handling.OHRCentralHandler)                  //office hour request stuff
	router.POST("/auth/login", handling.LoginHandler)                        //login authentication
	router.POST("/admin/:action", handling.POSTAdminCentralHandler)          //POST admin requests (rolelist, etc.)
	router.POST("/profile/:action", h.POSTProfileCentralHandler)             //POST profile requests

	handling.CacheClean()

	//tells the server to start listening to requests
	log.SetOutput(os.Stdout)
	log.Println("Go server running on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
