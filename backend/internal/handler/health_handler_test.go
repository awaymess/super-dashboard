package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthHandler_Health(t *testing.T) {
	gin.SetMode(gin.TestMode)

	healthHandler := NewHealthHandler()
	router := gin.New()
	healthHandler.RegisterHealthRoutes(router)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}
}

func TestHealthHandler_Live(t *testing.T) {
	gin.SetMode(gin.TestMode)

	healthHandler := NewHealthHandler()
	router := gin.New()
	healthHandler.RegisterHealthRoutes(router)

	req, err := http.NewRequest(http.MethodGet, "/health/live", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "alive" {
		t.Errorf("Expected status 'alive', got '%s'", response.Status)
	}
}

func TestHealthHandler_Ready(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		ready          bool
		checkers       []HealthChecker
		expectedStatus int
		expectedState  string
	}{
		{
			name:           "ready with no checkers",
			ready:          true,
			checkers:       nil,
			expectedStatus: http.StatusOK,
			expectedState:  "ready",
		},
		{
			name:  "ready with healthy checker",
			ready: true,
			checkers: []HealthChecker{
				func() (string, bool, string) {
					return "test", true, "healthy"
				},
			},
			expectedStatus: http.StatusOK,
			expectedState:  "ready",
		},
		{
			name:  "not ready with unhealthy checker",
			ready: true,
			checkers: []HealthChecker{
				func() (string, bool, string) {
					return "test", false, "unhealthy"
				},
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedState:  "not_ready",
		},
		{
			name:           "not ready state",
			ready:          false,
			checkers:       nil,
			expectedStatus: http.StatusServiceUnavailable,
			expectedState:  "not_ready",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			healthHandler := NewHealthHandler()
			healthHandler.SetReady(tt.ready)
			for _, checker := range tt.checkers {
				healthHandler.AddHealthChecker(checker)
			}

			router := gin.New()
			healthHandler.RegisterHealthRoutes(router)

			req, err := http.NewRequest(http.MethodGet, "/health/ready", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			var response HealthResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Status != tt.expectedState {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedState, response.Status)
			}
		})
	}
}

func TestHealthHandler_AddHealthChecker(t *testing.T) {
	healthHandler := NewHealthHandler()

	checker := func() (string, bool, string) {
		return "database", true, "connected"
	}

	healthHandler.AddHealthChecker(checker)

	if len(healthHandler.checkers) != 1 {
		t.Errorf("Expected 1 checker, got %d", len(healthHandler.checkers))
	}
}

func TestHealthHandler_SetReady(t *testing.T) {
	healthHandler := NewHealthHandler()

	// Default should be ready
	if !healthHandler.ready {
		t.Error("Expected default ready state to be true")
	}

	healthHandler.SetReady(false)
	if healthHandler.ready {
		t.Error("Expected ready state to be false after SetReady(false)")
	}

	healthHandler.SetReady(true)
	if !healthHandler.ready {
		t.Error("Expected ready state to be true after SetReady(true)")
	}
}
