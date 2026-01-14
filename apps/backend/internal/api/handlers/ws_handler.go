package handlers

import (
	"ai-task-manager/backend/internal/websocket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gorilla_websocket "github.com/gorilla/websocket"
)

// WSHandler handles WebSocket connections
type WSHandler struct {
	hub      *websocket.Hub
	upgrader gorilla_websocket.Upgrader
}

// NewWSHandler creates a new WSHandler
func NewWSHandler(hub *websocket.Hub) *WSHandler {
	return &WSHandler{
		hub: hub,
		upgrader: gorilla_websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development
				// In production, you should check the origin
				return true
			},
		},
	}
}

// HandleWebSocket handles WebSocket connection upgrades
func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	client := websocket.NewClient(h.hub, conn)
	h.hub.RegisterClient(client)

	// Start read and write pumps in separate goroutines
	go client.WritePump()
	go client.ReadPump()
}
