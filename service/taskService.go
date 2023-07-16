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

func (ts *TaskService) CreateTask(task types.Task, username string) error {
	taskdao := types.TaskDAO{
		ID: task.ID,
		Title: task.Title,
		Description: task.Description,
		Isdone: task.Isdone,
		Username: username,
	}
	err := ts.taskStorage.AddTask(taskdao)
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
func (ts *TaskService) FindAll() ([]types.Task, error) {
	tasks, err := ts.taskStorage.FindAll()
	if err != nil {
		log.Printf("Service, Error while retreving tasks ", err)
		return nil, err
	}
	return tasks, nil
}
func (ts *TaskService) FindTaskByUsername(username string) (*types.Task, error) {
	taskDAO, err := ts.taskStorage.FindTaskByUsername(username)
	if err != nil {
		log.Printf("Could not find Task by id: %v", err)
		return &types.Task{}, err
	}
	task := types.Task{
		ID: taskDAO.ID,
		Title: taskDAO.Title,
		Description: taskDAO.Description,
		Isdone: taskDAO.Isdone,
	}
	return &task, nil
}

func (ts *TaskService) FindAllByUserName(username string) ([]types.Task, error) {

	tasks, err := ts.taskStorage.FindAllByUsername(username)
	tasksDto := make([]types.Task, len(tasks))

	// Map list of type B to type A
	for i, b := range tasks {
		tasksDto[i] = types.Task{
			ID:   b.ID,
			Title: b.Title,
			Description: b.Description,
			Isdone: b.Isdone,
			}
	}
	if err != nil {
		log.Printf("Service, Error while retreving tasks ", err)
		return nil, err
	}
	return tasksDto, nil
}

func (ts *TaskService) DeleteById(id primitive.ObjectID) error {
	err := ts.taskStorage.DeleteTaskById(id)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil

}
func (ts *TaskService) UpdateTaskById(id primitive.ObjectID, task types.Task) error {
	err := ts.taskStorage.UpdateTaskById(id, task)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	return nil

}
