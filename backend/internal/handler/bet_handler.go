package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValueBetResponse represents a value bet response.
type ValueBetResponse struct {
	MatchID        string  `json:"match_id"`
	BetType        string  `json:"bet_type"`
	BookmakerOdds  float64 `json:"bookmaker_odds"`
	FairOdds       float64 `json:"fair_odds"`
	Value          float64 `json:"value"`
	Confidence     int     `json:"confidence"`
	KellyStake     float64 `json:"kelly_stake"`
	ExpectedValue  float64 `json:"expected_value"`
	Recommendation string  `json:"recommendation"`
}

// BettingStatsResponse represents betting statistics.
type BettingStatsResponse struct {
	TotalBets      int     `json:"total_bets"`
	WonBets        int     `json:"won_bets"`
	LostBets       int     `json:"lost_bets"`
	VoidBets       int     `json:"void_bets"`
	TotalStaked    float64 `json:"total_staked"`
	TotalReturns   float64 `json:"total_returns"`
	Profit         float64 `json:"profit"`
	ROI            float64 `json:"roi"`
	AvgOdds        float64 `json:"avg_odds"`
	AvgStake       float64 `json:"avg_stake"`
	WinRate        float64 `json:"win_rate"`
	CurrentStreak  int     `json:"current_streak"`
	BestStreak     int     `json:"best_streak"`
	WorstStreak    int     `json:"worst_streak"`
}

// BetRequest represents a bet placement request.
type BetRequest struct {
	MatchID string  `json:"match_id" binding:"required"`
	BetType string  `json:"bet_type" binding:"required"`
	Odds    float64 `json:"odds" binding:"required,gt=1"`
	Stake   float64 `json:"stake" binding:"required,gt=0"`
}

// BetHandler handles betting-related HTTP requests.
type BetHandler struct{}

// NewBetHandler creates a new BetHandler instance.
func NewBetHandler() *BetHandler {
	return &BetHandler{}
}

// GetValueBets returns value bets.
// @Summary Get value bets
// @Description Get a list of current value bets
// @Tags betting
// @Produce json
// @Success 200 {array} ValueBetResponse
// @Router /api/v1/betting/value-bets [get]
func (h *BetHandler) GetValueBets(c *gin.Context) {
	// Mock value bets data
	valueBets := []ValueBetResponse{
		{
			MatchID:        "1",
			BetType:        "Over 2.5 Goals",
			BookmakerOdds:  1.75,
			FairOdds:       1.60,
			Value:          9.4,
			Confidence:     72,
			KellyStake:     4.2,
			ExpectedValue:  8.5,
			Recommendation: "bet",
		},
		{
			MatchID:        "2",
			BetType:        "Home Win",
			BookmakerOdds:  1.55,
			FairOdds:       1.42,
			Value:          9.2,
			Confidence:     78,
			KellyStake:     5.1,
			ExpectedValue:  9.0,
			Recommendation: "bet",
		},
		{
			MatchID:        "3",
			BetType:        "BTTS Yes",
			BookmakerOdds:  1.85,
			FairOdds:       1.70,
			Value:          8.8,
			Confidence:     65,
			KellyStake:     3.5,
			ExpectedValue:  7.2,
			Recommendation: "bet",
		},
	}

	c.JSON(http.StatusOK, valueBets)
}

// GetBettingStats returns betting statistics.
// @Summary Get betting statistics
// @Description Get user's betting statistics
// @Tags betting
// @Produce json
// @Success 200 {object} BettingStatsResponse
// @Router /api/v1/betting/stats [get]
func (h *BetHandler) GetBettingStats(c *gin.Context) {
	stats := BettingStatsResponse{
		TotalBets:     156,
		WonBets:       105,
		LostBets:      48,
		VoidBets:      3,
		TotalStaked:   15600,
		TotalReturns:  18720,
		Profit:        3120,
		ROI:           20.0,
		AvgOdds:       1.85,
		AvgStake:      100,
		WinRate:       68.6,
		CurrentStreak: 3,
		BestStreak:    12,
		WorstStreak:   -5,
	}

	c.JSON(http.StatusOK, stats)
}

// PlaceBet handles bet placement.
// @Summary Place a bet
// @Description Place a new bet
// @Tags betting
// @Accept json
// @Produce json
// @Param request body BetRequest true "Bet details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/betting/bets [post]
func (h *BetHandler) PlaceBet(c *gin.Context) {
	var req BetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Calculate potential win
	potentialWin := req.Stake * req.Odds

	c.JSON(http.StatusCreated, gin.H{
		"id":            "bet_" + req.MatchID,
		"match_id":      req.MatchID,
		"bet_type":      req.BetType,
		"odds":          req.Odds,
		"stake":         req.Stake,
		"potential_win": potentialWin,
		"status":        "pending",
		"placed_at":     "2024-12-04T00:00:00Z",
	})
}

// GetBets returns user's bet history.
// @Summary Get bet history
// @Description Get user's betting history
// @Tags betting
// @Produce json
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/betting/bets [get]
func (h *BetHandler) GetBets(c *gin.Context) {
	// Mock bet history
	bets := []map[string]interface{}{
		{
			"id":            "bet_1",
			"match_id":      "1",
			"bet_type":      "Home Win",
			"odds":          2.10,
			"stake":         100,
			"potential_win": 210,
			"status":        "won",
			"profit":        110,
			"placed_at":     "2024-12-01T10:00:00Z",
			"settled_at":    "2024-12-01T12:00:00Z",
		},
		{
			"id":            "bet_2",
			"match_id":      "2",
			"bet_type":      "Over 2.5",
			"odds":          1.85,
			"stake":         50,
			"potential_win": 92.5,
			"status":        "lost",
			"profit":        -50,
			"placed_at":     "2024-12-02T14:00:00Z",
			"settled_at":    "2024-12-02T16:00:00Z",
		},
		{
			"id":            "bet_3",
			"match_id":      "3",
			"bet_type":      "BTTS Yes",
			"odds":          1.75,
			"stake":         75,
			"potential_win": 131.25,
			"status":        "pending",
			"placed_at":     "2024-12-04T09:00:00Z",
		},
	}

	c.JSON(http.StatusOK, bets)
}

// RegisterBetRoutes registers betting routes.
func (h *BetHandler) RegisterBetRoutes(rg *gin.RouterGroup) {
	betting := rg.Group("/betting")
	{
		betting.GET("/value-bets", h.GetValueBets)
		betting.GET("/stats", h.GetBettingStats)
		betting.GET("/bets", h.GetBets)
		betting.POST("/bets", h.PlaceBet)
	}
}
