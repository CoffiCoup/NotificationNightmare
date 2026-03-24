package handlers

//student view focused handlers here

import (
	"fmt"
	"log"
	"net/http"
)

func studentViewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Student view %s!", r.URL.Path[1:])
}

func student