package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/awaymess/super-dashboard/backend/internal/repository"
)

// MatchHandler handles match-related HTTP requests.
type MatchHandler struct {
	matchRepo repository.MatchRepository
}

// NewMatchHandler creates a new MatchHandler instance.
func NewMatchHandler(matchRepo repository.MatchRepository) *MatchHandler {
	return &MatchHandler{matchRepo: matchRepo}
}

// ListMatches returns all matches.
// @Summary List all matches
// @Description Get a list of all matches
// @Tags betting
// @Produce json
// @Success 200 {array} model.Match
// @Router /api/v1/betting/matches [get]
func (h *MatchHandler) ListMatches(c *gin.Context) {
	matches, err := h.matchRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch matches"})
		return
	}
	c.JSON(http.StatusOK, matches)
}

// GetMatch returns a single match by ID.
// @Summary Get match by ID
// @Description Get details of a specific match
// @Tags betting
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} model.Match
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/betting/matches/{id} [get]
func (h *MatchHandler) GetMatch(c *gin.Context) {
	id := c.Param("id")
	match, err := h.matchRepo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "match not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch match"})
		return
	}
	c.JSON(http.StatusOK, match)
}

// GetMatchOdds returns odds for a specific match.
// @Summary Get match odds
// @Description Get betting odds for a specific match
// @Tags betting
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {array} model.Odds
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/betting/matches/{id}/odds [get]
func (h *MatchHandler) GetMatchOdds(c *gin.Context) {
	id := c.Param("id")
	odds, err := h.matchRepo.GetOddsByMatchID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch odds"})
		return
	}
	c.JSON(http.StatusOK, odds)
}

// RegisterMatchRoutes registers match-related routes.
func (h *MatchHandler) RegisterMatchRoutes(rg *gin.RouterGroup) {
	betting := rg.Group("/betting")
	{
		betting.GET("/matches", h.ListMatches)
		betting.GET("/matches/:id", h.GetMatch)
		betting.GET("/matches/:id/odds", h.GetMatchOdds)
	}
}
