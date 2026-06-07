package main

import (
	"encoding/json"
	"net/http"

	"task-api-comparison/store"

	"github.com/gorilla/mux"
)

func listTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.All())
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	t, ok := store.ByID(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid body"})
		return
	}
	t := store.Create(input.Title, input.Description)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var input struct {
		Done bool `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid body"})
		return
	}
	t, ok := store.UpdateDone(id, input.Done)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/tasks", listTasks).Methods("GET")
	r.HandleFunc("/v1/tasks/{id}", getTask).Methods("GET")
	r.HandleFunc("/v1/tasks", createTask).Methods("POST")
	r.HandleFunc("/v1/tasks/{id}", updateTask).Methods("PATCH")

	http.ListenAndServe(":8082", r)
}
