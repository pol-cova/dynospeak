package handlers

import (
	"backend/pkg/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// Gorilla WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Global RoomManager instance
var roomManager = NewRoomManager()

func HandleWebSocket(c echo.Context, dbInstance *mongo.Collection) error {
	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	// Extract username and room name
	username := c.QueryParam("username")
	roomName := c.QueryParam("room_name")

	// Add WebSocket connection to the room
	roomManager.JoinRoom(roomName, ws)
	defer roomManager.LeaveRoom(roomName, ws)

	// Infinite loop to continuously listen for messages
	for {
		var msg models.Message
		// Read the message from WebSocket
		err := ws.ReadJSON(&msg)
		if err != nil {
			// If there's an error, break the loop
			break
		}

		// Add metadata to the message
		msg.ID = primitive.NewObjectID()
		msg.Username = username
		msg.RoomName = roomName
		msg.CreatedAt = time.Now()

		// Store the message in MongoDB
		_, err = dbInstance.InsertOne(c.Request().Context(), msg)
		if err != nil {
			// Log or handle the error as needed
			break
		}

		// Broadcast the message to all connected clients in the room
		roomManager.Broadcast(roomName, msg)
	}
	return nil
}
func CreateNewRoom(c echo.Context, dbInstance *mongo.Collection) error {
	// Extract the username from the context (assuming it's set by middleware)
	username := c.Get("username").(string)

	// Bind the request payload to the ChatRoom struct
	chatRoom := new(models.ChatRoom)
	if err := c.Bind(chatRoom); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request payload",
		})
	}

	// Set the owner of the chat room to the logged-in user
	chatRoom.Owner = username

	// Call the Create method to insert the chat room into the database
	if err := chatRoom.Create(dbInstance); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Could not create chat room: " + err.Error(),
		})
	}

	// Return the newly created chat room
	return c.JSON(http.StatusCreated, chatRoom)
}

func GetMessagesInRoom(c echo.Context, dbInstance *mongo.Collection) error {
	roomName := c.QueryParam("room_name")
	if roomName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Room name is required",
		})
	}

	// Define a filter to find messages in the specified room
	filter := bson.M{"room_name": roomName}
	cursor, err := dbInstance.Find(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to retrieve messages",
		})
	}
	defer cursor.Close(c.Request().Context())

	var messages []models.Message
	if err := cursor.All(c.Request().Context(), &messages); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to decode messages",
		})
	}

	return c.JSON(http.StatusOK, messages)
}

func GetAllRooms(c echo.Context, dbInstance *mongo.Collection) error {
	cursor, err := dbInstance.Find(c.Request().Context(), bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to retrieve chat rooms",
		})
	}
	defer cursor.Close(c.Request().Context())

	var rooms []models.ChatRoom
	if err := cursor.All(c.Request().Context(), &rooms); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to decode chat rooms",
		})
	}

	return c.JSON(http.StatusOK, rooms)
}
