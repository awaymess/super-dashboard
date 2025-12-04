package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("NewHub returned nil")
	}
	if hub.clients == nil {
		t.Error("Expected clients map to be initialized")
	}
	if hub.broadcast == nil {
		t.Error("Expected broadcast channel to be initialized")
	}
	if hub.register == nil {
		t.Error("Expected register channel to be initialized")
	}
	if hub.unregister == nil {
		t.Error("Expected unregister channel to be initialized")
	}
}

func TestHub_ClientCount(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	if hub.ClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", hub.ClientCount())
	}
}

func TestHub_RegisterUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create a mock client
	client := &Client{
		ID:            "test-client-1",
		hub:           hub,
		send:          make(chan []byte, 256),
		Subscriptions: make(map[string]bool),
	}

	// Register client
	hub.Register(client)
	time.Sleep(50 * time.Millisecond) // Give the hub time to process

	if hub.ClientCount() != 1 {
		t.Errorf("Expected 1 client after registration, got %d", hub.ClientCount())
	}

	// Unregister client
	hub.Unregister(client)
	time.Sleep(50 * time.Millisecond) // Give the hub time to process

	if hub.ClientCount() != 0 {
		t.Errorf("Expected 0 clients after unregistration, got %d", hub.ClientCount())
	}
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create mock clients
	client1 := &Client{
		ID:            "test-client-1",
		hub:           hub,
		send:          make(chan []byte, 256),
		Subscriptions: make(map[string]bool),
	}
	client2 := &Client{
		ID:            "test-client-2",
		hub:           hub,
		send:          make(chan []byte, 256),
		Subscriptions: make(map[string]bool),
	}

	hub.Register(client1)
	hub.Register(client2)
	time.Sleep(50 * time.Millisecond)

	// Broadcast an event
	event := Event{
		Type:    EventStockPriceUpdate,
		Payload: map[string]interface{}{"symbol": "AAPL", "price": 150.00},
	}

	err := hub.Broadcast(event)
	if err != nil {
		t.Errorf("Broadcast failed: %v", err)
	}

	// Both clients should receive the message
	timeout := time.After(500 * time.Millisecond)
	received1 := false
	received2 := false

	for !received1 || !received2 {
		select {
		case <-client1.send:
			received1 = true
		case <-client2.send:
			received2 = true
		case <-timeout:
			t.Fatal("Timeout waiting for broadcast messages")
		}
	}
}

func TestHub_BroadcastToChannel(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Create mock clients
	client1 := &Client{
		ID:            "test-client-1",
		hub:           hub,
		send:          make(chan []byte, 256),
		Subscriptions: make(map[string]bool),
	}
	client2 := &Client{
		ID:            "test-client-2",
		hub:           hub,
		send:          make(chan []byte, 256),
		Subscriptions: make(map[string]bool),
	}

	// Subscribe client1 to a channel
	client1.Subscribe("stocks")

	hub.Register(client1)
	hub.Register(client2)
	time.Sleep(50 * time.Millisecond)

	// Broadcast to the stocks channel
	event := Event{
		Type:    EventStockPriceUpdate,
		Payload: map[string]interface{}{"symbol": "AAPL", "price": 150.00},
	}

	err := hub.BroadcastToChannel("stocks", event)
	if err != nil {
		t.Errorf("BroadcastToChannel failed: %v", err)
	}

	// Only client1 should receive the message
	select {
	case <-client1.send:
		// Expected
	case <-time.After(500 * time.Millisecond):
		t.Error("Client1 should have received the message")
	}

	// Client2 should not receive the message
	select {
	case <-client2.send:
		t.Error("Client2 should not have received the message (not subscribed)")
	case <-time.After(100 * time.Millisecond):
		// Expected - no message received
	}
}

func TestClient_SubscribeUnsubscribe(t *testing.T) {
	client := &Client{
		ID:            "test-client",
		Subscriptions: make(map[string]bool),
	}

	// Test Subscribe
	client.Subscribe("channel1")
	if !client.Subscriptions["channel1"] {
		t.Error("Expected client to be subscribed to channel1")
	}

	client.Subscribe("channel2")
	if !client.Subscriptions["channel2"] {
		t.Error("Expected client to be subscribed to channel2")
	}

	// Test Unsubscribe
	client.Unsubscribe("channel1")
	if client.Subscriptions["channel1"] {
		t.Error("Expected client to be unsubscribed from channel1")
	}
	if !client.Subscriptions["channel2"] {
		t.Error("Expected client to still be subscribed to channel2")
	}
}

func TestEvent_Serialization(t *testing.T) {
	event := Event{
		Type:      EventMatchLiveScore,
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{"homeScore": 2, "awayScore": 1},
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	var decoded Event
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	if decoded.Type != event.Type {
		t.Errorf("Expected type %s, got %s", event.Type, decoded.Type)
	}
}

func TestWebSocketHandler_StatusEndpoint(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	handler := NewWebSocketHandler(hub)
	router := gin.New()
	handler.RegisterRoutes(router)

	// Test /ws/status endpoint
	req, _ := http.NewRequest(http.MethodGet, "/ws/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["connected_clients"]; !ok {
		t.Error("Expected connected_clients in response")
	}
	if _, ok := response["status"]; !ok {
		t.Error("Expected status in response")
	}
}

func TestWebSocketHandler_HandleWebSocket(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	handler := NewWebSocketHandler(hub)
	router := gin.New()
	handler.RegisterRoutes(router)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect as WebSocket client
	dialer := websocket.Dialer{}
	conn, resp, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("Expected status %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}

	// Wait for the hub to register the client
	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("Expected 1 connected client, got %d", hub.ClientCount())
	}
}

func TestWebSocketHandler_BroadcastToConnectedClients(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	handler := NewWebSocketHandler(hub)
	router := gin.New()
	handler.RegisterRoutes(router)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect as WebSocket client
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Wait for connection to be established
	time.Sleep(100 * time.Millisecond)

	// Broadcast an event
	event := Event{
		Type:    EventNotificationNew,
		Payload: map[string]interface{}{"message": "Hello, World!"},
	}
	_ = hub.Broadcast(event)

	// Read message from WebSocket
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	var received Event
	if err := json.Unmarshal(message, &received); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if received.Type != EventNotificationNew {
		t.Errorf("Expected event type %s, got %s", EventNotificationNew, received.Type)
	}
}

func TestHub_ConcurrentOperations(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	var wg sync.WaitGroup
	numClients := 10

	// Concurrently register clients
	clients := make([]*Client, numClients)
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			client := &Client{
				ID:            "test-client-" + string(rune('0'+idx)),
				hub:           hub,
				send:          make(chan []byte, 256),
				Subscriptions: make(map[string]bool),
			}
			clients[idx] = client
			hub.Register(client)
		}(i)
	}
	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != numClients {
		t.Errorf("Expected %d clients, got %d", numClients, hub.ClientCount())
	}

	// Concurrently broadcast messages
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			event := Event{
				Type:    EventStockPriceUpdate,
				Payload: map[string]interface{}{"index": idx},
			}
			_ = hub.Broadcast(event)
		}(i)
	}
	wg.Wait()

	// Concurrently unregister clients
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			hub.Unregister(clients[idx])
		}(i)
	}
	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("Expected 0 clients after unregistration, got %d", hub.ClientCount())
	}
}

func TestEventTypes(t *testing.T) {
	// Verify all event types are defined correctly
	eventTypes := []EventType{
		EventMatchLiveScore,
		EventMatchOddsUpdate,
		EventMatchStatusChange,
		EventBetResult,
		EventStockPriceUpdate,
		EventStockAlertTriggered,
		EventStockNews,
		EventNotificationNew,
		EventUserSessionExpired,
	}

	for _, et := range eventTypes {
		if string(et) == "" {
			t.Errorf("Event type has empty string value")
		}
	}
}
