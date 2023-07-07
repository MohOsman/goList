package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var tasks []Task
var currentTaskID = 1 // Variable to track the current task ID

func main() {
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/tasks", handlePostRequest).Methods("POST")
	myRouter.HandleFunc("/tasks/{id}", handleGetRequest).Methods("GET")
	myRouter.HandleFunc("/tasks", handleGetTasksRequest).Methods("GET")
	myRouter.HandleFunc("/tasks/{id}", handleUpdateRequest).Methods("PUT")
	myRouter.HandleFunc("/tasks/{id}", handleDeleteRequest).Methods("DELETE")

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(myRouter)))

}

// create read update and delete tasks
type Task struct {
	ID          int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Isdone      bool   `json:isdone`
}

func handlePostRequest(rw http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error deoading payload: %w", err)
		return
	}
	task.ID = currentTaskID
	currentTaskID++
	tasks = append(tasks, task)

	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "Post request Handled successfully")

}

func handleGetRequest(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"] // map with key id
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Invalid task ID: %v", err)
		return
	}

	var foundtTask *Task
	for _, task := range tasks {
		if task.ID == taskID {
			foundtTask = &task
			break
		}
	}

	jsonData, err := json.Marshal(foundtTask)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error marshaling task: %v", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonData)
	fmt.Print("Get request Handled successfully")
}

func handleGetTasksRequest(rw http.ResponseWriter, r *http.Request) {

	jsondata, err := json.Marshal(tasks)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Error marshaling task: %v", err)
		return
	}
	if IsEmptyList(tasks) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsondata)
	fmt.Print("Get request Handled  successfully")
}

func handleUpdateRequest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Invalid task ID: %v", err)
		return
	}
	var task Task
	err1 := json.NewDecoder(r.Body).Decode(&task)
	if err1 != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Error deoading payload: %w", err)
		return
	}

	var value uint

	for index, _ := range tasks {
		if task.ID == taskID {
			value = uint(index)
			break
		}
		tasks[value] = task
		rw.WriteHeader(http.StatusCreated)
		fmt.Fprintf(rw, "Put request Handled successfully")

	}

}

func handleDeleteRequest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "Invalid task ID: %v", err)
		return
	}
	var index = -1
	for i, task := range tasks {
		if task.ID == taskID {
			index = i
			break
		}
	}
	if index != 1 {
		tasks = append(tasks[:index], tasks[index+1:]...)

	}

	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Delete request Handled successfully")

}

func IsEmptyList(tasks []Task) bool {
	return len(tasks) == 0
}
