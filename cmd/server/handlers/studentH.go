package handlers

//student view focused handlers here

import (
	"fmt"
	"net/http"
)

func StudentViewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Student view %s!", r.URL.Path[1:])
}

func StudentMakeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Student make %s!", r.URL.Path[1:])
}

func StudentUpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Student update %s!", r.URL.Path[1:])
}
