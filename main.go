package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var (
	todos  = []Todo{}
	nextID = 1
	mu     sync.Mutex
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", getTodos)
	mux.HandleFunc("POST /todos", createTodo)

	http.ListenAndServe(":8080", enableCORS(mux))
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	if newTodo.Title == "" {
		http.Error(w, "Título obrigatório", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newTodo.ID = nextID
	nextID++
	newTodo.Completed = false
	todos = append(todos, newTodo)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
