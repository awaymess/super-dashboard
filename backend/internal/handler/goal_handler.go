package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"super-dashboard/backend/internal/model"
	"super-dashboard/backend/internal/repository"
)

// GoalHandler handles goal-related HTTP requests.
type GoalHandler struct {
	goalRepo *repository.GoalRepository
}

// NewGoalHandler creates a new GoalHandler.
func NewGoalHandler(goalRepo *repository.GoalRepository) *GoalHandler {
	return &GoalHandler{
		goalRepo: goalRepo,
	}
}

// CreateGoal handles POST /api/goals
func (h *GoalHandler) CreateGoal(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		GoalType    string  `json:"goal_type" binding:"required"`
		Description string  `json:"description" binding:"required"`
		TargetValue float64 `json:"target_value" binding:"required"`
		Timeframe   string  `json:"timeframe" binding:"required"`
		TargetDate  string  `json:"target_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetDate, err := time.Parse("2006-01-02", req.TargetDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid target date format"})
		return
	}

	goal := &model.Goal{
		UserID:       userID.(uuid.UUID),
		GoalType:     req.GoalType,
		Description:  req.Description,
		TargetValue:  req.TargetValue,
		CurrentValue: 0,
		Progress:     0,
		Timeframe:    req.Timeframe,
		TargetDate:   targetDate,
		Achieved:     false,
	}

	if err := h.goalRepo.CreateGoal(c.Request.Context(), goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"goal": goal})
}

// GetGoals handles GET /api/goals
func (h *GoalHandler) GetGoals(c *gin.Context) {
	userID, _ := c.Get("user_id")

	goals, err := h.goalRepo.GetUserGoals(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goals": goals})
}

// GetActiveGoals handles GET /api/goals/active
func (h *GoalHandler) GetActiveGoals(c *gin.Context) {
	userID, _ := c.Get("user_id")

	goals, err := h.goalRepo.GetActiveGoals(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goals": goals})
}

// GetGoalByID handles GET /api/goals/:id
func (h *GoalHandler) GetGoalByID(c *gin.Context) {
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goal ID"})
		return
	}

	goal, err := h.goalRepo.GetGoalByID(c.Request.Context(), goalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goal": goal})
}

// UpdateGoal handles PUT /api/goals/:id
func (h *GoalHandler) UpdateGoal(c *gin.Context) {
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goal ID"})
		return
	}

	goal, err := h.goalRepo.GetGoalByID(c.Request.Context(), goalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}

	var req struct {
		Description  *string  `json:"description"`
		TargetValue  *float64 `json:"target_value"`
		CurrentValue *float64 `json:"current_value"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Description != nil {
		goal.Description = *req.Description
	}
	if req.TargetValue != nil {
		goal.TargetValue = *req.TargetValue
	}
	if req.CurrentValue != nil {
		if err := h.goalRepo.UpdateGoalProgress(c.Request.Context(), goalID, *req.CurrentValue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Reload goal
		goal, _ = h.goalRepo.GetGoalByID(c.Request.Context(), goalID)
		c.JSON(http.StatusOK, gin.H{"goal": goal})
		return
	}

	if err := h.goalRepo.UpdateGoal(c.Request.Context(), goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goal": goal})
}

// MarkGoalAchieved handles PUT /api/goals/:id/achieved
func (h *GoalHandler) MarkGoalAchieved(c *gin.Context) {
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goal ID"})
		return
	}

	if err := h.goalRepo.MarkGoalAsAchieved(c.Request.Context(), goalID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal marked as achieved"})
}

// DeleteGoal handles DELETE /api/goals/:id
func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goal ID"})
		return
	}

	if err := h.goalRepo.DeleteGoal(c.Request.Context(), goalID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}

// GetGoalStatistics handles GET /api/goals/statistics
func (h *GoalHandler) GetGoalStatistics(c *gin.Context) {
	userID, _ := c.Get("user_id")

	stats, err := h.goalRepo.GetGoalStatistics(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
