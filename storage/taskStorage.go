package storage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"goList/types"
)

type TaskStorage interface {
	AddTask(task types.Task) error
	FindTaskById(primitive.ObjectID) (*types.Task, error)
	

}
