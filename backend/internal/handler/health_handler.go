package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints.
type HealthHandler struct {
	mu       sync.RWMutex
	ready    bool
	checkers []HealthChecker
}

// HealthChecker is a function that checks a dependency's health.
type HealthChecker func() (name string, healthy bool, message string)

// HealthResponse represents a health check response.
type HealthResponse struct {
	Status  string                 `json:"status"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// DependencyStatus represents the status of a single dependency.
type DependencyStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		ready:    true,
		checkers: make([]HealthChecker, 0),
	}
}

// AddHealthChecker adds a health checker function.
func (h *HealthHandler) AddHealthChecker(checker HealthChecker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers = append(h.checkers, checker)
}

// SetReady sets the readiness state.
func (h *HealthHandler) SetReady(ready bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.ready = ready
}

// Health returns basic health status.
// @Summary Basic health check
// @Description Returns basic health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status: "ok",
	})
}

// Ready checks if the service is ready to accept traffic.
// @Summary Readiness check
// @Description Checks if the service and all dependencies are ready
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	h.mu.RLock()
	ready := h.ready
	checkers := h.checkers
	h.mu.RUnlock()

	if !ready {
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status: "not_ready",
		})
		return
	}

	details := make(map[string]interface{})
	allHealthy := true

	for _, checker := range checkers {
		name, healthy, message := checker()
		status := "up"
		if !healthy {
			status = "down"
			allHealthy = false
		}
		details[name] = DependencyStatus{
			Status:  status,
			Message: message,
		}
	}

	if !allHealthy {
		c.JSON(http.StatusServiceUnavailable, HealthResponse{
			Status:  "not_ready",
			Details: details,
		})
		return
	}

	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ready",
		Details: details,
	})
}

// Live checks if the service is alive.
// @Summary Liveness check
// @Description Checks if the service is alive (used by Kubernetes liveness probe)
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status: "alive",
	})
}

// RegisterHealthRoutes registers health check routes.
func (h *HealthHandler) RegisterHealthRoutes(r *gin.Engine) {
	r.GET("/health", h.Health)
	r.GET("/health/ready", h.Ready)
	r.GET("/health/live", h.Live)
}
