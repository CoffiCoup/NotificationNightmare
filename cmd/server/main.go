package main

//the starting point for the project for everything

import (
	"fmt"
	"log"
	"net/http"
	//"html/template"
)

// example handler
func examplehandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "booltest %s!", r.URL.Path[1:])
}

// run handler through this when we need to check for certain certification before request can be handled
func certificationRequired(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//rudimentary check currently, we need to add in specifics
		if r.TLS != nil {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Client certificate required", http.StatusUnauthorized)
		}
	}
}

func main() {
	//when making a new handler, make sure to add it here accordingly
	http.HandleFunc("/", examplehandler)
	http.HandleFunc("/secure", certificationRequired(http.HandlerFunc(examplehandler)))

	//certificate and key file paths, make sure the files are actually here
	certFile := "cert.crt"
	keyFile := "private.key"

	err := http.ListenAndServeTLS(":8080", certFile, keyFile, nil)
	if err != nil {
		log.Fatalf("HTTPS server failed: %v", err)
	}
}
