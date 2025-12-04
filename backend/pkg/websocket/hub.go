// Package websocket provides WebSocket hub for real-time event broadcasting.
package websocket

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// EventType represents the type of WebSocket event.
type EventType string

const (
	// Sports events
	EventMatchLiveScore   EventType = "match:live_score"
	EventMatchOddsUpdate  EventType = "match:odds_update"
	EventMatchStatusChange EventType = "match:status_change"
	EventBetResult        EventType = "bet:result"

	// Stock events
	EventStockPriceUpdate   EventType = "stock:price_update"
	EventStockAlertTriggered EventType = "stock:alert_triggered"
	EventStockNews          EventType = "stock:news"

	// System events
	EventNotificationNew    EventType = "notification:new"
	EventUserSessionExpired EventType = "user:session_expired"
)

// Event represents a WebSocket event.
type Event struct {
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

// Client represents a WebSocket client connection.
type Client struct {
	ID            string
	UserID        string
	Subscriptions map[string]bool // Channels the client is subscribed to
	Messages      chan []byte
	mu            sync.RWMutex
}

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Info().Str("client_id", client.ID).Msg("WebSocket client connected")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Messages)
			}
			h.mu.Unlock()
			log.Info().Str("client_id", client.ID).Msg("WebSocket client disconnected")

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Messages <- message:
				default:
					// Client buffer full, remove it
					h.mu.RUnlock()
					h.mu.Lock()
					close(client.Messages)
					delete(h.clients, client)
					h.mu.Unlock()
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends an event to all connected clients.
func (h *Hub) Broadcast(event Event) error {
	event.Timestamp = time.Now()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	h.broadcast <- data
	return nil
}

// BroadcastToChannel sends an event to clients subscribed to a specific channel.
func (h *Hub) BroadcastToChannel(channel string, event Event) error {
	event.Timestamp = time.Now()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		client.mu.RLock()
		if client.Subscriptions[channel] {
			select {
			case client.Messages <- data:
			default:
				// Skip if buffer full
			}
		}
		client.mu.RUnlock()
	}
	return nil
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Subscribe adds a client to a channel.
func (c *Client) Subscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Subscriptions == nil {
		c.Subscriptions = make(map[string]bool)
	}
	c.Subscriptions[channel] = true
}

// Unsubscribe removes a client from a channel.
func (c *Client) Unsubscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Subscriptions, channel)
}

// WebSocketHandler handles WebSocket upgrade requests.
// Note: This is a stub implementation using Server-Sent Events (SSE) for simplicity.
// For production, use gorilla/websocket or similar library.
type WebSocketHandler struct {
	hub *Hub
}

// NewWebSocketHandler creates a new WebSocket handler.
func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HandleSSE handles Server-Sent Events connections as a WebSocket alternative.
// This is a simpler implementation that doesn't require additional dependencies.
func (h *WebSocketHandler) HandleSSE(c *gin.Context) {
	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Create client
	client := &Client{
		ID:            c.GetHeader("X-Request-ID"),
		Messages:      make(chan []byte, 256),
		Subscriptions: make(map[string]bool),
	}

	// Register client
	h.hub.register <- client

	// Ensure unregister on disconnect
	defer func() {
		h.hub.unregister <- client
	}()

	// Send events to client using SSE format
	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-client.Messages:
			if !ok {
				return false
			}
			c.SSEvent("message", string(msg))
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

// RegisterRoutes registers WebSocket routes.
func (h *WebSocketHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ws", h.HandleSSE)
	r.GET("/events", h.HandleSSE) // Alias for SSE

	// Status endpoint
	r.GET("/ws/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"connected_clients": h.hub.ClientCount(),
			"status":           "operational",
		})
	})
}
