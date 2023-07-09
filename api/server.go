package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"goList/types"
	"goList/utils"
	"log"
	"net/http"
	"strconv"
)
type Server struct {
	listenAddr string
}

var tasks []types.Task
var currentTaskID = 1 // Variable to track the current task ID


func NewServer(lisntAddr string) *Server{
	return &Server{listenAddr: lisntAddr}
}

func (s *Server) Start() error  {
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


	return http.ListenAndServe(":8080", corsMiddleware(myRouter))

}


func handlePostRequest(rw http.ResponseWriter, r *http.Request) {
	var task types.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println("Error deoading payload: %w", err)
		return
	}
	task.ID = currentTaskID
	currentTaskID++
	tasks = append(tasks, task)
	createdTask, err := json.Marshal(task)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(rw, "Error marshaling task: %v", err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	log.Println( "Post request Handled successfully")
	rw.Write(createdTask)
}

func handleGetRequest(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"] // map with key id
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Invalid task ID: %v", err)
		return
	}

	var foundtTask *types.Task
	for _, task := range tasks {
		if task.ID == taskID {
			foundtTask = &task
			break
		}
	}

	jsonData, err := json.Marshal(foundtTask)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(rw, "Error marshaling task: %v", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	log.Println("Get request Handled successfully")
	rw.Write(jsonData)
}

func handleGetTasksRequest(rw http.ResponseWriter, r *http.Request) {

	jsondata, err := json.Marshal(tasks)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(rw, "Error marshaling task: %v", err)
		return
	}
	if utils.IsEmptyList(tasks) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	log.Println("Get request handled successfully ")
	rw.Write(jsondata)
}

func handleUpdateRequest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Invalid task ID: %v", err)
		return
	}
	var task types.Task
	err1 := json.NewDecoder(r.Body).Decode(&task)
	if err1 != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Error deoading payload: %w", err)
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
		log.Println("Put request Handled successfully")

	}

}

func handleDeleteRequest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Invalid task ID: %v", err)
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
	currentTaskID = 0;

	rw.WriteHeader(http.StatusOK)
}

