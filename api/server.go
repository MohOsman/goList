package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goList/service"
	"goList/types"
	"goList/utils"

	"log"
	"net/http"
	"strconv"
)

type Server struct {
	listenAddr   string
	userService1 service.UserService
	ts           service.TaskService
}

var tasks []types.Task

var currentTaskID = 1 // Variable to track the current task ID

func NewServer(lisntAddr string,
	userService service.UserService,
	taskService service.TaskService) *Server {
	return &Server{listenAddr: lisntAddr,
		userService1: userService,
		ts:           taskService}
}

func (s *Server) Start() error {
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/tasks", s.handlePostTask).Methods("POST")
	myRouter.HandleFunc("/tasks/{id}", s.handleGetRequest).Methods("GET")
	myRouter.HandleFunc("/tasks", handleGetTasksRequest).Methods("GET")
	myRouter.HandleFunc("/tasks/{id}", handleUpdateRequest).Methods("PUT")
	myRouter.HandleFunc("/tasks/{id}", handleDeleteRequest).Methods("DELETE")
	myRouter.HandleFunc("/register", s.handleRegisterUser).Methods("POST")

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

func (s *Server) handleRegisterUser(rw http.ResponseWriter, r *http.Request) {
	var newuser types.User
	err := json.NewDecoder(r.Body).Decode(&newuser)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println("Error deoading payload: %w", err)
		return
	}
	user := s.userService1.RegisterUser(newuser)
	createdUser, err := json.Marshal(user)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(rw, "Error marshaling task: %v", err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	log.Println("user created Handled successfully")
	rw.Write(createdUser)

}

func (s *Server) handlePostTask(rw http.ResponseWriter, r *http.Request) {
	var task types.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println("Error deoading payload: %w", err)
		return
	}

	err1 := s.ts.CreateTask(task)
	if err1 != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error saviong the database %w", err)
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	log.Println("Post request Handled successfully")

}

func (s *Server) handleGetRequest(rw http.ResponseWriter, r *http.Request) {

	taskID := mux.Vars(r)["id"]
	// map with key id
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Invalid task ID: %v", err)
		return
	}

	foundTask, err := s.ts.FindTaskById(objectID)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		log.Println("Could not find task with id") // look how you can id  inside the log message
	}
	jsonData, err := json.Marshal(foundTask)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(rw, "Error marshaling task: %v", err)
		return
	}
	log.Println(taskID)

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
	log.Print(taskID)

	var value uint
	for index, _ := range tasks {
		value = uint(index)
		break
	}
	tasks[value] = task
	rw.WriteHeader(http.StatusCreated)
	log.Println("Put request Handled successfully")

}

func handleDeleteRequest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	test, err := strconv.Atoi(idStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Invalid task ID: %v", err)
		return
	}
	log.Print(test)

	rw.WriteHeader(http.StatusOK)
}
