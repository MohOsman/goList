package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *UserRepository) RegisterUser(user types.UserDAO) error {

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	// Create the index
	_, err := r.collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Printf("Error creating index %v", err)
		return err
	}
	_, err = r.collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Error storing user %v", err)
		return err
	}
	return nil

}

func (r *UserRepository) FindUserByUsername(username string) (*types.UserDAO, error) {
	filter := bson.M{"username": username}
	var user types.UserDAO
	err := r.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Printf("Error while retrving user form database %v", err)
		return nil, err
	}
	return &user, err
}
