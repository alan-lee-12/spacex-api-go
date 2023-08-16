package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Define a handler function for the GET endpoint
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// Write a response to the client
		fmt.Fprintf(w, "Hello, world!")
	})

	// Start the HTTP server on port 8080
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
