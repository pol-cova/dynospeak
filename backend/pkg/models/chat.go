package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (c *ChatRoom) Create(db *mongo.Collection) error {
	// Check if the room name already exists
	filter := bson.M{"room_name": c.RoomName}
	count, err := db.CountDocuments(context.Background(), filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("chat room with this name already exists")
	}

	// Generate a new ObjectID for the chat room
	c.ID = primitive.NewObjectID()

	// Insert the new chat room into the MongoDB collection
	_, err = db.InsertOne(context.Background(), c)
	if err != nil {
		return err
	}

	return nil
}
