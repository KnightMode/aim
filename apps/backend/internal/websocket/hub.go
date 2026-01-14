package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// RegisterClient registers a new client with the hub
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient unregisters a client from the hub
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// Message types for WebSocket events
type WSMessage struct {
	Type      string      `json:"type"`
	TaskID    uint        `json:"task_id,omitempty"`
	Status    string      `json:"status,omitempty"`
	LogLevel  string      `json:"log_level,omitempty"`
	Message   string      `json:"message,omitempty"`
	Result    string      `json:"result,omitempty"`
	Error     string      `json:"error,omitempty"`
	AgentName string      `json:"agent_name,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// BroadcastTaskStatusChanged broadcasts a task status change event
func (h *Hub) BroadcastTaskStatusChanged(taskID uint, status string, timestamp string) {
	msg := WSMessage{
		Type:      "task_status_changed",
		TaskID:    taskID,
		Status:    status,
		Timestamp: timestamp,
	}
	h.broadcastMessage(msg)
}

// BroadcastExecutionLog broadcasts an execution log event
func (h *Hub) BroadcastExecutionLog(taskID uint, logLevel string, message string, timestamp string) {
	msg := WSMessage{
		Type:      "execution_log",
		TaskID:    taskID,
		LogLevel:  logLevel,
		Message:   message,
		Timestamp: timestamp,
	}
	h.broadcastMessage(msg)
}

// BroadcastTaskCompleted broadcasts a task completion event
func (h *Hub) BroadcastTaskCompleted(taskID uint, result string, timestamp string) {
	msg := WSMessage{
		Type:      "task_completed",
		TaskID:    taskID,
		Result:    result,
		Timestamp: timestamp,
	}
	h.broadcastMessage(msg)
}

// BroadcastTaskFailed broadcasts a task failure event
func (h *Hub) BroadcastTaskFailed(taskID uint, errorMsg string, timestamp string) {
	msg := WSMessage{
		Type:      "task_failed",
		TaskID:    taskID,
		Error:     errorMsg,
		Timestamp: timestamp,
	}
	h.broadcastMessage(msg)
}

// BroadcastAgentStarted broadcasts an agent started event
func (h *Hub) BroadcastAgentStarted(taskID uint, agentName string, timestamp string) {
	msg := WSMessage{
		Type:      "agent_started",
		TaskID:    taskID,
		AgentName: agentName,
		Timestamp: timestamp,
	}
	h.broadcastMessage(msg)
}

// broadcastMessage is a helper to marshal and broadcast a message
func (h *Hub) broadcastMessage(msg WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}
	h.Broadcast(data)
}
