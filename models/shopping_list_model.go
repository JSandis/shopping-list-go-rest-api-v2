package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type List struct {
	ID      primitive.ObjectID `bson:"_id"`
	User_id string             `json:"user_id"`
}

type ListItem struct {
	ID      primitive.ObjectID `bson:"_id"`
	Name    string             `json:"name" validate:"required"`
	Status  bool               `json:"status"`
	List_id string             `json:"list_id"`
}
