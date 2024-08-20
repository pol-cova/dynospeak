package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatRoom represents a chat room in the application.
type ChatRoom struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Owner    string             `json:"owner" bson:"owner" binding:"required"`
	RoomName string             `json:"room_name" bson:"room_name" binding:"required"`
}

// Message represents a single message in a chat room.
type Message struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username" binding:"required"`
	RoomName  string             `json:"room_name" bson:"room_name" binding:"required"`
	Message   string             `json:"message" bson:"message" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
