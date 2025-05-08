package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	SystemResponse string `json:"system-response"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Чтение данных из тела запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Декодирование JSON в структуру
	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Формирование ответа
	response := Response{
		SystemResponse: fmt.Sprintf("Hello, %s!", req.Name),
	}

	// Устанавливаем заголовок Content-Type как JSON
	w.Header().Set("Content-Type", "application/json")

	// Отправляем ответ в формате JSON
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error sending response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)

	// Запуск сервера на порту 8080
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

