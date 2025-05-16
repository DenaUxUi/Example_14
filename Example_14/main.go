package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Response struct {
    Message string `json:"system-response"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    // Разрешаем CORS
    w.Header().Set("Access-Control-Allow-Origin", "*") // или конкретный домен вместо '*'
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Предварительный запрос OPTIONS
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != http.MethodGet {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    name := r.URL.Query().Get("name")
    if name == "" {
        http.Error(w, "Missing 'name' parameter in query", http.StatusBadRequest)
        return
    }

    response := Response{
        Message: fmt.Sprintf("Hello, %s!", name),
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

