package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goList/authentication"
	"goList/service"
	"goList/types"
	"goList/utils"

	"log"
	"net/http"
)

type Server struct {
	listenAddr   string
	userService1 service.UserService
	ts           service.TaskService
}

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
	myRouter.HandleFunc("/tasks", s.handleGetTasksRequest).Methods("GET")
	myRouter.HandleFunc("/tasks/{id}", s.handleUpdateRequest).Methods("PUT")
	myRouter.HandleFunc("/tasks/{id}", s.handleDeleteRequest).Methods("DELETE")
	myRouter.HandleFunc("/register", s.handleRegisterUser).Methods("POST")
	myRouter.HandleFunc("/login", s.handlelogin).Methods("POST")

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	return http.ListenAndServe(":8080", corsMiddleware(myRouter))

}

// users
// my own router for users ?

func (s *Server) handleRegisterUser(rw http.ResponseWriter, r *http.Request) {
	var newuser types.User
	err := json.NewDecoder(r.Body).Decode(&newuser)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println("Error deoading payload: %w", err)
		return
	}
	err = s.userService1.RegisterUser(newuser)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println("Username already exsit: %w", err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	log.Println("user created handled successfully")

}

func (s *Server) handlelogin(rw http.ResponseWriter, r *http.Request) {
	var user types.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println("Error decoding payload: %w", err)
		return
	}
	username, err := s.userService1.Login(user)

	if err != nil {
		rw.WriteHeader(http.StatusForbidden)
		log.Println("Wrong credentials", err)
		return
	}

	token, err := authentication.GenerateJWT(*username)
	if err != nil {
		log.Println("Error gernerating token: %w", err)
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"token":"` + token + `"}`))
}

// Tasks
func (s *Server) handlePostTask(rw http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		log.Println("Authorization Token in not present in the request")
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, err := authentication.ValidateAndExtractUsername(token)
	if err != nil {
		log.Printf("could not procces token %v", err)
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return

	}
	var task types.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println("Error deoading payload: %w", err)
		return
	}

	err1 := s.ts.CreateTask(task, username)
	if err1 != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error saviong the database %w", err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")
	log.Println("Post request handled successfully")

}

func (s *Server) handleGetRequest(rw http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username, err := authentication.ValidateAndExtractUsername(token)
	if err != nil {
		log.Printf("could not procces token %v", err)
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return

	}
	id := mux.Vars(r)["id"]

	taskID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Invalid task ID: %v", err)
		return
	}
	foundTask, err := s.ts.FindTaskByUsername(username, taskID)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		log.Printf("Could not find task with username %v", username)
		return
	}

	jsonData, err := json.Marshal(foundTask)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(rw, "Error marshaling task: %v", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	log.Println("Get All tasks request handled successfully")
	rw.Write(jsonData)
}

func (s *Server) handleGetTasksRequest(rw http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username, err := authentication.ValidateAndExtractUsername(token)
	if err != nil {
		log.Printf("could not procces token %v", err)
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return

	}

	tasks, err := s.ts.FindAllByUserName(username)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		log.Println(rw, "Error retreving tasks: %v", err)
		return

	}
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
	log.Println("Get All request handled successfully ")
	_, err = rw.Write(jsondata)
	if err != nil {
		log.Print("Error http Writing")
		return
	}
}

func (s *Server) handleUpdateRequest(rw http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username, err := authentication.ValidateAndExtractUsername(token)

	if err != nil {
		log.Printf("could not procces token %v", err)
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return

	}
	id := mux.Vars(r)["id"]

	taskID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Invalid task ID: %v", err)
		return
	}
	var task types.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Error deoading payload: %w", err)
		return
	}
	err = s.ts.UpdateTaskByUsername(username, taskID, task)
	if err != nil {
		log.Printf("Error while updating task %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	log.Println("Put request handled successfully")

}

func (s *Server) handleDeleteRequest(rw http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username, err := authentication.ValidateAndExtractUsername(token)

	if err != nil {
		log.Printf("could not procces token %v", err)
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return

	}
	Id := mux.Vars(r)["id"]
	taskId, err := primitive.ObjectIDFromHex(Id)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(rw, "Invalid task ID: %v", err)
		return
	}
	err = s.ts.DeleteByUsername(username, taskId)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Println(rw, "Invalid task ID: %v", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
