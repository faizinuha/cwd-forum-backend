package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-quickstart/config"
	"gin-quickstart/internal/model"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WSHub struct {
	clients map[uint]*websocket.Conn
	mutex   sync.RWMutex
}

func NewWSHub() *WSHub {
	return &WSHub{
		clients: make(map[uint]*websocket.Conn),
	}
}

func (h *WSHub) Register(userID uint, conn *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	// Close existing connection if any to prevent memory leaks from dangling connections
	if existingConn, ok := h.clients[userID]; ok {
		existingConn.Close()
	}
	h.clients[userID] = conn
}

func (h *WSHub) Unregister(userID uint) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if conn, ok := h.clients[userID]; ok {
		conn.Close()
		delete(h.clients, userID)
	}
}

// StartRedisListener runs in a background goroutine, subscribing to notifications
func (h *WSHub) StartRedisListener() {
	if config.RedisClient == nil {
		log.Println("RedisClient is nil, cannot start WS Redis Listener")
		return
	}

	ctx := context.Background()
	pubsub := config.RedisClient.Subscribe(ctx, "realtime:notifications")
	defer pubsub.Close()

	ch := pubsub.Channel()
	log.Println("WebSocket Redis Listener started. Listening on 'realtime:notifications'")

	for msg := range ch {
		var notification model.Notification
		err := json.Unmarshal([]byte(msg.Payload), &notification)
		if err != nil {
			log.Println("Error unmarshaling notification:", err)
			continue
		}

		// Find if the user is currently connected to this hub instance
		h.mutex.RLock()
		conn, ok := h.clients[notification.UserId]
		h.mutex.RUnlock()

		if ok {
			// Send notification to the connected client
			err = conn.WriteJSON(notification)
			if err != nil {
				log.Println("Error writing JSON to websocket for user", notification.UserId, ":", err)
				h.Unregister(notification.UserId)
			}
		}
	}
}
