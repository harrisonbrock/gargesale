package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	h := http.HandlerFunc(Echo)
	log.Println("Listening on localhost:8000")
	if err := http.ListenAndServe("localhost:8000", h); err != nil {
		log.Fatal(err)
	}
}

// Echo is a basic HTTP Handler.
func Echo(w http.ResponseWriter, r * http.Request) {
	fmt.Fprintln(w, "You ask to", r.Method, r.URL.Path)
}