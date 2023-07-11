package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goList/storage"
	"goList/types"
	"log"
)

type TaskService struct {
	taskStorage storage.TaskStorage
}

func NewTaskService(storage storage.TaskStorage) *TaskService {
	return &TaskService{
		taskStorage: storage,
	}

}

func (ts *TaskService) CreateTask(task types.Task) error {
	err := ts.taskStorage.AddTask(task)
	if err != nil {
		log.Printf("service, Could not insert the task into the database")
		return err
	}
	return nil
}

func (ts *TaskService) FindTaskById(id primitive.ObjectID) (*types.Task, error) {
	task, err := ts.taskStorage.FindTaskById(id)
	if err != nil {
		log.Printf("Could not find Task by id: %v", err)
		return &types.Task{}, err
	}
	return task, nil
}
