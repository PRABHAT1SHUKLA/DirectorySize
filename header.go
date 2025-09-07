package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// The `http.ResponseWriter` (w) implements the `io.Writer` interface.
	// fmt.Fprintf can use it to write a formatted string directly to the HTTP response body.
	name := "Alice"
	fmt.Fprintf(w, "Hello, %s! You requested kkdslkdsldk: %s\n", name, r.URL.Path)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	// Get the value of the User-Agent header
	userAgent := r.Header.Get("User-Agent")
	fmt.Fprintf(w, userAgent)
}

func bitch(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "xjckxmckxmkcmxkcmkxmc")
}

func rem(w http.ResponseWriter, r *http.Request) {
	// --- Receiving Headers ---
	// The http.Request (r) has a Header field, which is a map[string][]string.
	// You can use the convenience method Get() to retrieve the first value for a header.
	requestedBy := r.Header.Get("X-Requested-By")
	log.Printf("Received request from: %s", requestedBy)

	// You can also iterate through all the received headers
	log.Println("Received headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("  %s: %s", name, value)
		}
	}

	// --- Setting Headers ---
	// The http.ResponseWriter (w) also has a Header() method that returns a Header map.
	// You must set headers before writing the response body.
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Custom-Header", "Golang-is-awesome")
	w.Header().Set("sdsds", "Golaesome")
	// Write the response body
	fmt.Fprintf(w, "Hello! Your request was received.\n")
	fmt.Fprintf(w, "We read the 'X-Requested-By' header as: %s\n", requestedBy)
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/bit", bitch)
	http.HandleFunc("/r", rem)
	log.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
