package handler

import (
	"github.com/gin-gonic/gin"
)

// Handler holds all HTTP handlers
type Handler struct {
	// TODO: Add service dependencies
}

// NewHandler creates a new Handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// TODO: Register routes
}
