package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"super-dashboard/backend/internal/model"
)

// GoalRepository handles database operations for user goals.
type GoalRepository struct {
	db *gorm.DB
}

// NewGoalRepository creates a new GoalRepository.
func NewGoalRepository(db *gorm.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

// CreateGoal creates a new goal.
func (r *GoalRepository) CreateGoal(ctx context.Context, goal *model.Goal) error {
	return r.db.WithContext(ctx).Create(goal).Error
}

// GetGoalByID retrieves a goal by ID.
func (r *GoalRepository) GetGoalByID(ctx context.Context, id uuid.UUID) (*model.Goal, error) {
	var goal model.Goal
	err := r.db.WithContext(ctx).First(&goal, id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// GetUserGoals retrieves all goals for a user.
func (r *GoalRepository) GetUserGoals(ctx context.Context, userID uuid.UUID) ([]model.Goal, error) {
	var goals []model.Goal
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&goals).Error
	return goals, err
}

// GetActiveGoals retrieves active (not completed) goals for a user.
func (r *GoalRepository) GetActiveGoals(ctx context.Context, userID uuid.UUID) ([]model.Goal, error) {
	var goals []model.Goal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND achieved = false", userID).
		Order("target_date ASC").
		Find(&goals).Error
	return goals, err
}

// GetGoalsByType retrieves goals filtered by type.
func (r *GoalRepository) GetGoalsByType(ctx context.Context, userID uuid.UUID, goalType string) ([]model.Goal, error) {
	var goals []model.Goal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND goal_type = ?", userID, goalType).
		Order("created_at DESC").
		Find(&goals).Error
	return goals, err
}

// GetGoalsByTimeframe retrieves goals filtered by timeframe.
func (r *GoalRepository) GetGoalsByTimeframe(ctx context.Context, userID uuid.UUID, timeframe string) ([]model.Goal, error) {
	var goals []model.Goal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND timeframe = ?", userID, timeframe).
		Order("created_at DESC").
		Find(&goals).Error
	return goals, err
}

// UpdateGoal updates a goal.
func (r *GoalRepository) UpdateGoal(ctx context.Context, goal *model.Goal) error {
	return r.db.WithContext(ctx).Save(goal).Error
}

// UpdateGoalProgress updates the current progress of a goal.
func (r *GoalRepository) UpdateGoalProgress(ctx context.Context, goalID uuid.UUID, currentValue float64) error {
	goal, err := r.GetGoalByID(ctx, goalID)
	if err != nil {
		return err
	}

	goal.CurrentValue = currentValue
	
	// Calculate progress percentage
	if goal.TargetValue > 0 {
		goal.Progress = (currentValue / goal.TargetValue) * 100
		if goal.Progress > 100 {
			goal.Progress = 100
		}
	}

	// Check if goal is achieved
	if currentValue >= goal.TargetValue && !goal.Achieved {
		goal.Achieved = true
		now := time.Now()
		goal.AchievedAt = &now
	}

	return r.UpdateGoal(ctx, goal)
}

// MarkGoalAsAchieved marks a goal as achieved.
func (r *GoalRepository) MarkGoalAsAchieved(ctx context.Context, goalID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.Goal{}).
		Where("id = ?", goalID).
		Updates(map[string]interface{}{
			"achieved":    true,
			"achieved_at": now,
			"progress":    100.0,
		}).Error
}

// DeleteGoal deletes a goal.
func (r *GoalRepository) DeleteGoal(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Goal{}, id).Error
}

// GetOverdueGoals retrieves goals that are past their target date but not achieved.
func (r *GoalRepository) GetOverdueGoals(ctx context.Context, userID uuid.UUID) ([]model.Goal, error) {
	var goals []model.Goal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND achieved = false AND target_date < ?", userID, time.Now()).
		Order("target_date ASC").
		Find(&goals).Error
	return goals, err
}

// GetUpcomingGoals retrieves goals with target dates in the near future.
func (r *GoalRepository) GetUpcomingGoals(ctx context.Context, userID uuid.UUID, days int) ([]model.Goal, error) {
	endDate := time.Now().AddDate(0, 0, days)
	var goals []model.Goal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND achieved = false AND target_date BETWEEN ? AND ?", 
			userID, time.Now(), endDate).
		Order("target_date ASC").
		Find(&goals).Error
	return goals, err
}

// GetGoalStatistics calculates goal statistics for a user.
func (r *GoalRepository) GetGoalStatistics(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	var stats struct {
		TotalGoals    int
		ActiveGoals   int
		AchievedGoals int
		OverdueGoals  int
		AvgProgress   float64
	}

	now := time.Now()

	err := r.db.WithContext(ctx).
		Model(&model.Goal{}).
		Where("user_id = ?", userID).
		Select(`
			COUNT(*) as total_goals,
			SUM(CASE WHEN achieved = false THEN 1 ELSE 0 END) as active_goals,
			SUM(CASE WHEN achieved = true THEN 1 ELSE 0 END) as achieved_goals,
			SUM(CASE WHEN achieved = false AND target_date < ? THEN 1 ELSE 0 END) as overdue_goals,
			AVG(CASE WHEN achieved = false THEN progress ELSE NULL END) as avg_progress
		`, now).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	achievementRate := 0.0
	if stats.TotalGoals > 0 {
		achievementRate = float64(stats.AchievedGoals) / float64(stats.TotalGoals) * 100
	}

	return map[string]interface{}{
		"total_goals":      stats.TotalGoals,
		"active_goals":     stats.ActiveGoals,
		"achieved_goals":   stats.AchievedGoals,
		"overdue_goals":    stats.OverdueGoals,
		"avg_progress":     stats.AvgProgress,
		"achievement_rate": achievementRate,
	}, nil
}

// GetGoalsByPriority retrieves goals sorted by progress and target date.
func (r *GoalRepository) GetGoalsByPriority(ctx context.Context, userID uuid.UUID) ([]model.Goal, error) {
	var goals []model.Goal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND achieved = false", userID).
		Order("target_date ASC, progress ASC").
		Find(&goals).Error
	return goals, err
}
