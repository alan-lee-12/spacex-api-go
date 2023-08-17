package main

import (
	"fmt"
	"net/http"
)

func main() {
	// GET /hello
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// Reject if not GET
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Write a response to the client
		fmt.Fprintf(w, "Hello, world!")
	})

	// GET all rockets
	http.HandleFunc("/rockets", getRocketsHandler)

	// Start the HTTP server on port 8080
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
