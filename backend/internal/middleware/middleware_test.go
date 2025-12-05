package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/awaymess/super-dashboard/backend/internal/model"
	"github.com/awaymess/super-dashboard/backend/internal/service"
)

// mockAuthService is a mock implementation of AuthService for testing.
type mockAuthService struct {
	jwtSecret string
}

func newMockAuthService() *mockAuthService {
	return &mockAuthService{
		jwtSecret: "test-secret",
	}
}

func (m *mockAuthService) Register(email, password, name string) (*model.User, error) {
	return nil, nil
}

func (m *mockAuthService) Login(email, password string) (string, string, error) {
	return "", "", nil
}

func (m *mockAuthService) RefreshToken(refreshToken string) (string, error) {
	return "", nil
}

func (m *mockAuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, service.ErrInvalidToken
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, service.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, service.ErrInvalidToken
	}

	return &claims, nil
}

func (m *mockAuthService) generateToken(userID, email, role string) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(m.jwtSecret))
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockAuthService()
	userID := uuid.New().String()
	validToken := mockService.generateToken(userID, "test@example.com", "user")

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{
			name:       "valid token",
			authHeader: "Bearer " + validToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing auth header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid format - no Bearer",
			authHeader: validToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid format - wrong prefix",
			authHeader: "Basic " + validToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AuthMiddleware(mockService))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestOptionalAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := newMockAuthService()
	userID := uuid.New().String()
	validToken := mockService.generateToken(userID, "test@example.com", "user")

	tests := []struct {
		name          string
		authHeader    string
		wantStatus    int
		expectUserID  bool
	}{
		{
			name:         "valid token",
			authHeader:   "Bearer " + validToken,
			wantStatus:   http.StatusOK,
			expectUserID: true,
		},
		{
			name:         "missing auth header",
			authHeader:   "",
			wantStatus:   http.StatusOK,
			expectUserID: false,
		},
		{
			name:         "invalid token - continues without user",
			authHeader:   "Bearer invalid-token",
			wantStatus:   http.StatusOK,
			expectUserID: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(OptionalAuthMiddleware(mockService))
			router.GET("/test", func(c *gin.Context) {
				_, exists := c.Get("user_id")
				c.JSON(http.StatusOK, gin.H{"has_user": exists})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestRoleMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		requiredRoles []string
		userRole      string
		wantStatus    int
	}{
		{
			name:          "user has required role",
			requiredRoles: []string{"admin"},
			userRole:      "admin",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "user has one of required roles",
			requiredRoles: []string{"admin", "moderator"},
			userRole:      "moderator",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "user does not have required role",
			requiredRoles: []string{"admin"},
			userRole:      "user",
			wantStatus:    http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("role", tt.userRole)
				c.Next()
			})
			router.Use(RoleMiddleware(tt.requiredRoles...))
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestAdminMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		userRole   string
		wantStatus int
	}{
		{
			name:       "admin user",
			userRole:   "admin",
			wantStatus: http.StatusOK,
		},
		{
			name:       "regular user",
			userRole:   "user",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("role", tt.userRole)
				c.Next()
			})
			router.Use(AdminMiddleware())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestInMemoryRateLimiter(t *testing.T) {
	limiter := NewInMemoryRateLimiter(RateLimitConfig{
		Requests: 3,
		Window:   time.Minute,
	})

	ctx := context.Background()

	// First 3 requests should be allowed
	for i := 0; i < 3; i++ {
		allowed, remaining, _, err := limiter.Allow(ctx, "test-key")
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
		if remaining != 3-i-1 {
			t.Errorf("Expected remaining %d, got %d", 3-i-1, remaining)
		}
	}

	// 4th request should be denied
	allowed, remaining, _, err := limiter.Allow(ctx, "test-key")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if allowed {
		t.Error("4th request should be denied")
	}
	if remaining != 0 {
		t.Errorf("Expected remaining 0, got %d", remaining)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limiter := NewInMemoryRateLimiter(RateLimitConfig{
		Requests: 2,
		Window:   time.Minute,
	})

	router := gin.New()
	router.Use(RateLimitMiddleware(limiter, "test"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// First 2 requests should succeed
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status %d, got %d", i+1, http.StatusOK, w.Code)
		}
	}

	// 3rd request should be rate limited
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("3rd request: expected status %d, got %d", http.StatusTooManyRequests, w.Code)
	}

	// Check rate limit headers
	if w.Header().Get("X-RateLimit-Remaining") != "0" {
		t.Errorf("Expected X-RateLimit-Remaining to be 0")
	}
	if w.Header().Get("X-RateLimit-Limit") != "2" {
		t.Errorf("Expected X-RateLimit-Limit to be 2, got %s", w.Header().Get("X-RateLimit-Limit"))
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := DefaultSecurityHeadersConfig()
	config.EnableHSTS = true

	router := gin.New()
	router.Use(SecurityHeadersMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check CSP header
	csp := w.Header().Get("Content-Security-Policy")
	if csp == "" {
		t.Error("Expected Content-Security-Policy header to be set")
	}

	// Check HSTS header
	hsts := w.Header().Get("Strict-Transport-Security")
	if hsts == "" {
		t.Error("Expected Strict-Transport-Security header to be set")
	}

	// Check X-Frame-Options header
	xfo := w.Header().Get("X-Frame-Options")
	if xfo != "DENY" {
		t.Errorf("Expected X-Frame-Options to be 'DENY', got '%s'", xfo)
	}

	// Check X-Content-Type-Options header
	xcto := w.Header().Get("X-Content-Type-Options")
	if xcto != "nosniff" {
		t.Errorf("Expected X-Content-Type-Options to be 'nosniff', got '%s'", xcto)
	}

	// Check Referrer-Policy header
	rp := w.Header().Get("Referrer-Policy")
	if rp == "" {
		t.Error("Expected Referrer-Policy header to be set")
	}

	// Check Cache-Control header
	cc := w.Header().Get("Cache-Control")
	if cc == "" {
		t.Error("Expected Cache-Control header to be set")
	}
}

func TestSecurityHeadersMiddleware_NoHSTS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := DefaultSecurityHeadersConfig()
	config.EnableHSTS = false

	router := gin.New()
	router.Use(SecurityHeadersMiddleware(config))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// HSTS should not be set when disabled
	hsts := w.Header().Get("Strict-Transport-Security")
	if hsts != "" {
		t.Errorf("Expected no Strict-Transport-Security header, got '%s'", hsts)
	}
}
