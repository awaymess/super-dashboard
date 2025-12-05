package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, check origin properly
		return true
	},
}

// Handler handles WebSocket upgrade requests.
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler.
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// ServeWS handles WebSocket requests from clients.
func (h *Handler) ServeWS(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Create client
	client := NewClient(h.hub, conn, userID.(uint))
	h.hub.register <- client

	// Start client pumps
	go client.WritePump()
	go client.ReadPump()
}

// RegisterRoutes registers WebSocket routes.
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/ws", h.ServeWS)
}

// Broadcaster provides methods to broadcast updates.
type Broadcaster struct {
	hub *Hub
}

// NewBroadcaster creates a new broadcaster.
func NewBroadcaster(hub *Hub) *Broadcaster {
	return &Broadcaster{
		hub: hub,
	}
}

// BroadcastOddsUpdate broadcasts odds update to subscribers.
func (b *Broadcaster) BroadcastOddsUpdate(update OddsUpdate) error {
	msg := NewOutgoingMessage("odds_update", ChannelOdds, update)
	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	b.hub.BroadcastToChannel(ChannelOdds, data)
	return nil
}

// BroadcastStockPrice broadcasts stock price update.
func (b *Broadcaster) BroadcastStockPrice(update StockPriceUpdate) error {
	msg := NewOutgoingMessage("stock_update", ChannelStockPrices, update)
	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	b.hub.BroadcastToChannel(ChannelStockPrices, data)
	return nil
}

// BroadcastMatchUpdate broadcasts live match update.
func (b *Broadcaster) BroadcastMatchUpdate(update MatchUpdate) error {
	msg := NewOutgoingMessage("match_update", ChannelMatchUpdates, update)
	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	b.hub.BroadcastToChannel(ChannelMatchUpdates, data)
	return nil
}

// BroadcastPortfolioUpdate broadcasts portfolio value update.
func (b *Broadcaster) BroadcastPortfolioUpdate(update PortfolioUpdate) error {
	msg := NewOutgoingMessage("portfolio_update", ChannelPortfolio, update)
	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	b.hub.BroadcastToChannel(ChannelPortfolio, data)
	return nil
}

// BroadcastAlert broadcasts alert notification.
func (b *Broadcaster) BroadcastAlert(alert AlertNotification) error {
	msg := NewOutgoingMessage("alert", ChannelAlerts, alert)
	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	b.hub.BroadcastToChannel(ChannelAlerts, data)
	return nil
}

// BroadcastNews broadcasts news update.
func (b *Broadcaster) BroadcastNews(news NewsUpdate) error {
	msg := NewOutgoingMessage("news_update", ChannelNews, news)
	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	b.hub.BroadcastToChannel(ChannelNews, data)
	return nil
}

// Marshal serializes the outgoing message to JSON.
func (m *OutgoingMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
