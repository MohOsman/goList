package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Isdone      bool               `json:isdone`
}
