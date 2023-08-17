package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	http.HandleFunc("/rockets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			return
		}

		// Make a call to SpaceX API to get all rocket data.
		apiUrl := "https://api.spacexdata.com/v4/rockets"
		res, err := http.Get(apiUrl)
		if err != nil {
			http.Error(w, "Server error.", http.StatusInternalServerError)
			fmt.Println("original error:", err)
			return
		}

		/*
			From ChatGPT:
			HTTP responses in Go, including the one returned by http.Get(), provide a response body stream that needs
			to be closed when you're done using it. If you fail to close the response body, it could lead to resource
			leaks or potential issues like keeping connections open longer than necessary, which can cause problems
			in high-traffic scenarios.
			This line schedules the response.Body.Close() function call to be executed when the surrounding function
			(in this case, the /rockets handler function) returns. This ensures that the response body is closed
			regardless of how the function exits, whether it's due to successful processing or encountering an error.
		*/
		defer res.Body.Close()

		// Read all the bytes from the response body stream and returns the content as a byte slice ([]byte)
		data, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, "Error reading API response.", http.StatusInternalServerError)
			return
		}
		fmt.Println("API response status code:", res.StatusCode)
		// fmt.Println("raw response body:", string(data))

		// Watch out for non-2XX responses.
		if res.StatusCode > 299 {
			errStr := fmt.Sprintf("API error: %v", string(data))
			http.Error(w, errStr, http.StatusInternalServerError)
			return
		}

		// Unmarshall data to JSON so it's printable.
		var parsed interface{}
		// Here we unmarshal JSON into a pointer.
		// Unmarshal unmarshals the JSON into the value pointed at by the pointer. If the pointer is nil, Unmarshal
		// allocates a new value for it to point to.
		if err := json.Unmarshal(data, &parsed); err != nil {
			http.Error(w, "Error parsing API response", http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}
		// Prettify, indent the raw JSON.
		prettyJSON, err := json.MarshalIndent(parsed, "", "\t")
		if err != nil {
			http.Error(w, "Error formatting JSON.", http.StatusInternalServerError)
		}
		fmt.Println(string(prettyJSON))

		// Write the response from the API directly to the client's response writer.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.StatusCode)
		w.Write(data)
	})

	// Start the HTTP server on port 8080
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
