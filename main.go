package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var tasks []Task
var currentTaskID = 1 // Variable to track the current task ID

func main() {
	http.HandleFunc("/create", handlePostRequest)
	http.HandleFunc("/tasks/", handleGetRequest)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

// create read update and delete tasks
type Task struct {
	ID          int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Isdone      bool   `json:isdone`
}

func handlePostRequest(rw http.ResponseWriter, r *http.Request) {

	// Todo check method if it post 
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
	fmt.Fprintf(rw, "Pos request Handled successfully")

}

func handleGetRequest(rw http.ResponseWriter, r *http.Request) {
	// Todo check method if it is get
   idStr := r.URL.Path[len("/tasks/"):]
   taskID, err := strconv.Atoi(idStr)
  if err != nil {
	rw.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(rw, "Invalid task ID: %v", err)
	return
  }
  
  var foundtTask *Task
  for  _,task:= range tasks{
	if task.ID == taskID{
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


}
