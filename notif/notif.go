package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "booltest", r.URL.Path[1:])
}

func main() {
	http.HandlerFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
