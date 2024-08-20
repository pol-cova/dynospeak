package handlers

import (
	"backend/pkg/models"
	"github.com/gorilla/websocket"
	"sync"
)

type RoomManager struct {
	rooms map[string]map[*websocket.Conn]bool
	mu    sync.Mutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]map[*websocket.Conn]bool),
	}
}

func (rm *RoomManager) JoinRoom(roomName string, conn *websocket.Conn) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if _, ok := rm.rooms[roomName]; !ok {
		rm.rooms[roomName] = make(map[*websocket.Conn]bool)
	}
	rm.rooms[roomName][conn] = true
}

func (rm *RoomManager) LeaveRoom(roomName string, conn *websocket.Conn) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if clients, ok := rm.rooms[roomName]; ok {
		delete(clients, conn)
		if len(clients) == 0 {
			delete(rm.rooms, roomName)
		}
	}
}

func (rm *RoomManager) Broadcast(roomName string, msg models.Message) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	for conn := range rm.rooms[roomName] {
		if err := conn.WriteJSON(msg); err != nil {
			conn.Close()
			delete(rm.rooms[roomName], conn)
		}
	}
}
