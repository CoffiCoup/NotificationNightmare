package main

import (
	"log"
	"net/http"
	"notif/cmd/server/handling"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()
	router.GET("/view/:page", handling.AuthMiddleware(handling.ViewHandler)) //viewing
	router.POST("/oh/:action", handling.OHCentralHandler)                    //office hour stuff
	router.POST("/ohr/:action", handling.OHRCentralHandler)                  //office hour request stuff
	router.POST("/auth/login", handling.LoginHandler)                        //login authentication
	router.GET("/fetch/:file", handling.FetchCentralHandler)                 //sending files to html
	router.GET("/profile/:action/:extra", handling.GETProfileCentralHandler) //GET profile requests
	router.POST("/profile/:action", handling.POSTProfileCentralHandler)      //POST profile requests

	handling.CacheClean()

	//tells the server to start listening to requests
	log.SetOutput(os.Stdout)
	log.Println("Go server running on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
