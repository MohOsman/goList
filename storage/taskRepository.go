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

func (ts *TaskRepository) AddTask(task types.Task) error {
	_, err := ts.taskColletion.InsertOne(context.TODO(), task)
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

func (ts *TaskRepository) FindAll() ([]types.Task, error) {
	var tasks []types.Task
	cursor, err := ts.taskColletion.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Print("Error while retrieving tasks: ", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task types.Task
		if err := cursor.Decode(&task); err != nil {
			log.Printf("Error decoding task: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error with the cursor: %v", err)
		return nil, err
	}

	return tasks, nil
}

func (ts *TaskRepository) DeleteTaskById(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := ts.taskColletion.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Error deleting task,%v ", err)
		return err
	}
	return nil
}

func (ts *TaskRepository) UpdateTaskById(id primitive.ObjectID, task types.Task) error {
	filter := bson.M{"_id": id}
	_, err := ts.taskColletion.ReplaceOne(context.TODO(), filter, task)
	if err != nil {
		log.Printf("Error while replacing task %v", err)
		return err
	}
	return nil
}

func (ts *TaskRepository) FindTaskByUsername(username string, id primitive.ObjectID) (*types.Task, error) {
	filter := bson.M{
		"_id":      id,
		"username": username,
	}
	var task types.Task
	err := ts.taskColletion.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		log.Printf("Error while retrving user form database %v", err)
		return nil, err
	}
	return &task, err
}

func (ts *TaskRepository) FindAllByUsername(username string) ([]types.Task, error) {
	var tasks []types.Task
	cursor, err := ts.taskColletion.Find(context.TODO(), bson.M{"username": username})
	if err != nil {
		log.Print("Error while retrieving tasks: ", err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var task types.Task
		if err := cursor.Decode(&task); err != nil {
			log.Printf("Error decoding task: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error with the cursor: %v", err)
		return nil, err
	}

	return tasks, nil
}

func (ts *TaskRepository) DeleteTaskByUsername(username string, id primitive.ObjectID) error {
	task, err := ts.FindTaskByUsername(username, id)
	if err != nil {
		return err
	}

	err = ts.DeleteTaskById(task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TaskRepository) UpdateTaskByUsername(username string, id primitive.ObjectID, task types.Task) error {
	rTask, err := ts.FindTaskByUsername(username, id)
	if err != nil {
		return err
	}
	newTask := types.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Isdone:      task.Isdone,
		Username:    rTask.Username,
	}
	err = ts.UpdateTaskById(rTask.ID, newTask)
	if err != nil {
		return err
	}

	return nil
}
