package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskDAO struct {
	ID          primitive.ObjectID
	Title       string            
	Description string            
	Isdone      bool             
	Username    string
}