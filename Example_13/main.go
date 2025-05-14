// server.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"system-response"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		http.Error(w, "Invalid JSON or missing 'name'", http.StatusBadRequest)
		return
	}

	response := Response{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

