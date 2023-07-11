package storage

import (
	context "context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"goList/types"
	"log"
)

type TaskRepository struct {
	taskColletion *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) *TaskRepository {
	return &TaskRepository{taskColletion: collection}
}

func (t *TaskRepository) AddTask(task types.Task) error {
	_, err := t.taskColletion.InsertOne(context.TODO(), task)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}
	return nil
}

func (ts *TaskRepository) FindTaskById(id primitive.ObjectID) (*types.Task, error) {
	filter := bson.M{"_id": id}

	var task types.Task
	err := ts.taskColletion.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	return &task, nil
}

//DeleteTaskById(int) error
