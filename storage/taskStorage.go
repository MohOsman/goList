package storage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goList/types"
)

type TaskStorage interface {
	AddTask(task types.TaskDAO) error
	FindTaskById(id primitive.ObjectID) (*types.Task, error)
	FindTaskByUsername(username string) (*types.TaskDAO, error)
	FindAll() ([]types.Task, error)
	FindAllByUsername(username string) ([]types.TaskDAO, error)  
	DeleteTaskById(id primitive.ObjectID) error
	UpdateTaskById(id primitive.ObjectID, task types.Task) error
}
