// Package websocket provides WebSocket hub for real-time event broadcasting using gorilla/websocket.
package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// EventType represents the type of WebSocket event.
type EventType string

const (
	// Sports events
	EventMatchLiveScore    EventType = "match:live_score"
	EventMatchOddsUpdate   EventType = "match:odds_update"
	EventMatchStatusChange EventType = "match:status_change"
	EventBetResult         EventType = "bet:result"

	// Stock events
	EventStockPriceUpdate    EventType = "stock:price_update"
	EventStockAlertTriggered EventType = "stock:alert_triggered"
	EventStockNews           EventType = "stock:news"

	// System events
	EventNotificationNew    EventType = "notification:new"
	EventUserSessionExpired EventType = "user:session_expired"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development
		// In production, you should check the origin
		return true
	},
}

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
	conn          *websocket.Conn
	send          chan []byte
	hub           *Hub
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
		broadcast:  make(chan []byte, 256),
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
				close(client.send)
			}
			h.mu.Unlock()
			log.Info().Str("client_id", client.ID).Msg("WebSocket client disconnected")

		case message := <-h.broadcast:
			h.mu.RLock()
			// Collect clients that need to be removed
			var clientsToRemove []*Client
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client buffer full, mark for removal
					clientsToRemove = append(clientsToRemove, client)
				}
			}
			h.mu.RUnlock()

			// Remove clients with full buffers
			if len(clientsToRemove) > 0 {
				h.mu.Lock()
				for _, client := range clientsToRemove {
					if _, ok := h.clients[client]; ok {
						close(client.send)
						delete(h.clients, client)
					}
				}
				h.mu.Unlock()
			}
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
			case client.send <- data:
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

// Register registers a client with the hub.
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister unregisters a client from the hub.
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
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

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Str("client_id", c.ID).Msg("WebSocket read error")
			}
			break
		}
		// Handle incoming messages (e.g., subscription requests)
		c.handleMessage(message)
	}
}

// handleMessage processes incoming messages from clients.
func (c *Client) handleMessage(message []byte) {
	// Parse message and handle subscription requests
	var msg struct {
		Action  string `json:"action"`
		Channel string `json:"channel"`
	}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Warn().Err(err).Str("client_id", c.ID).Msg("Failed to parse client message")
		return
	}

	switch msg.Action {
	case "subscribe":
		if msg.Channel != "" {
			c.Subscribe(msg.Channel)
			log.Info().Str("client_id", c.ID).Str("channel", msg.Channel).Msg("Client subscribed to channel")
		}
	case "unsubscribe":
		if msg.Channel != "" {
			c.Unsubscribe(msg.Channel)
			log.Info().Str("client_id", c.ID).Str("channel", msg.Channel).Msg("Client unsubscribed from channel")
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WebSocketHandler handles WebSocket upgrade requests using gorilla/websocket.
type WebSocketHandler struct {
	hub *Hub
}

// NewWebSocketHandler creates a new WebSocket handler.
func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HandleWebSocket handles WebSocket connections.
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade connection to WebSocket")
		return
	}

	// Generate client ID
	clientID := c.GetHeader("X-Request-ID")
	if clientID == "" {
		clientID = uuid.New().String()
	}

	// Create client
	client := &Client{
		ID:            clientID,
		hub:           h.hub,
		conn:          conn,
		send:          make(chan []byte, 256),
		Subscriptions: make(map[string]bool),
	}

	// Register client
	h.hub.register <- client

	// Start read and write pumps
	go client.writePump()
	go client.readPump()
}

// RegisterRoutes registers WebSocket routes.
func (h *WebSocketHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ws", h.HandleWebSocket)

	// Status endpoint
	r.GET("/ws/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"connected_clients": h.hub.ClientCount(),
			"status":            "operational",
		})
	})
}
