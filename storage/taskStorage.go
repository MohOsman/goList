package storage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goList/types"
)

type TaskStorage interface {
	AddTask(task types.Task) error
	FindTaskById(id primitive.ObjectID) (*types.Task, error)
	FindAll() ([]types.Task, error)
	DeleteTaskById(id primitive.ObjectID) error
	UpdateTaskById(id primitive.ObjectID, task types.Task) error
}
