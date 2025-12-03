package handler

import (
	"github.com/gin-gonic/gin"
)

// Handler holds all HTTP handlers
type Handler struct {
	// TODO: Add service dependencies
}

// NewHandler creates a new handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// RegisterRoutes registers all routes
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// TODO: Register routes
}
