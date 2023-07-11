package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"goList/types"
	"log"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(Usercollection *mongo.Collection) *UserRepository {
	return &UserRepository{
		collection: Usercollection,
	}
}

func (r *UserRepository) RegisterUser(user types.User) error {
	_, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return err
	}
	return nil
}
