package main

import (
	"context"
	"goList/service"

	"log"
	"goList/api"
	"goList/storage"
)

func main() {
	client, err := storage.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Create the database and collection
	db, err := storage.CreateDatabaseAndCollections(client, "todolist")
	if err != nil {
		log.Fatal(err)
	}

	// beans
	userCollection := db.Collection("users")
	taskColletion := db.Collection("tasks")
	userRepository := storage.NewUserRepository(userCollection)
	taskRepository := storage.NewTaskRepository(taskColletion)
	userService := service.NewUserService(userRepository)
	taskService := service.NewTaskService(taskRepository)

	server := api.NewServer(":8080", *userService, *taskService)
	log.Fatal(server.Start())

}
