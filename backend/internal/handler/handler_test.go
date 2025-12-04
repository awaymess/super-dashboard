package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()

	// Register the health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Create a request to the health endpoint
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	expected := `{"status":"ok"}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, w.Body.String())
	}
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockAuthService()
	handler := NewAuthHandler(mockService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterAuthRoutes(v1)

	tests := []struct {
		name       string
		body       RegisterRequest
		wantStatus int
	}{
		{
			name: "valid registration",
			body: RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid email",
			body: RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
				Name:     "Test User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "short password",
			body: RegisterRequest{
				Email:    "test2@example.com",
				Password: "short",
				Name:     "Test User",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing name",
			body: RegisterRequest{
				Email:    "test3@example.com",
				Password: "password123",
				Name:     "",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockAuthService()
	handler := NewAuthHandler(mockService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterAuthRoutes(v1)

	// First register a user
	registerBody := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	bodyBytes, _ := json.Marshal(registerBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	tests := []struct {
		name       string
		body       LoginRequest
		wantStatus int
	}{
		{
			name: "valid login",
			body: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong password",
			body: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "non-existent user",
			body: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var response LoginResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.AccessToken == "" {
					t.Error("Expected access token to be set")
				}
				if response.RefreshToken == "" {
					t.Error("Expected refresh token to be set")
				}
			}
		})
	}
}

func TestAuthHandler_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockAuthService()
	handler := NewAuthHandler(mockService)

	router := gin.New()
	v1 := router.Group("/api/v1")
	handler.RegisterAuthRoutes(v1)

	// First register and login to get a refresh token
	registerBody := RegisterRequest{
		Email:    "refresh@example.com",
		Password: "password123",
		Name:     "Refresh User",
	}
	bodyBytes, _ := json.Marshal(registerBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	loginBody := LoginRequest{
		Email:    "refresh@example.com",
		Password: "password123",
	}
	bodyBytes, _ = json.Marshal(loginBody)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var loginResponse LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &loginResponse); err != nil {
		t.Fatalf("Failed to unmarshal login response: %v", err)
	}

	tests := []struct {
		name       string
		body       RefreshRequest
		wantStatus int
	}{
		{
			name: "valid refresh",
			body: RefreshRequest{
				RefreshToken: loginResponse.RefreshToken,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid refresh token",
			body: RefreshRequest{
				RefreshToken: "invalid-token",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "empty refresh token",
			body: RefreshRequest{
				RefreshToken: "",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var response RefreshResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.AccessToken == "" {
					t.Error("Expected access token to be set")
				}
			}
		})
	}
}
