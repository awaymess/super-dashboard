package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/awaymess/super-dashboard/backend/internal/model"
	"github.com/awaymess/super-dashboard/backend/internal/repository"
)

// findMockDataPathForHandler finds the mock data directory for handler tests
func findMockDataPathForHandler() string {
	// Try relative paths from different working directories
	paths := []string{
		"../../mock",
		"../../../mock",
		"mock",
	}
	for _, p := range paths {
		matchesPath := filepath.Join(p, "matches.json")
		if _, err := os.Stat(matchesPath); err == nil {
			return matchesPath
		}
	}
	return "../../mock/matches.json"
}

func setupMatchHandlerRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockPath := findMockDataPathForHandler()
	matchRepo, err := repository.NewMockMatchRepository(mockPath)
	if err != nil {
		t.Fatalf("Failed to create mock match repository: %v", err)
	}

	handler := NewMatchHandler(matchRepo)
	v1 := router.Group("/api/v1")
	handler.RegisterMatchRoutes(v1)

	return router
}

func TestMatchHandler_ListMatches(t *testing.T) {
	router := setupMatchHandlerRouter(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/betting/matches", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var matches []model.Match
	if err := json.Unmarshal(w.Body.Bytes(), &matches); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Based on mock data, we expect 5 matches
	if len(matches) != 5 {
		t.Errorf("Expected 5 matches, got %d", len(matches))
	}

	// Verify each match has required fields
	for _, match := range matches {
		if match.ID.String() == "" || match.ID.String() == "00000000-0000-0000-0000-000000000000" {
			t.Error("Match ID should be a valid UUID")
		}
		if match.League == "" {
			t.Error("Match League should not be empty")
		}
		if match.HomeTeam.Name == "" {
			t.Error("Match HomeTeam.Name should not be empty")
		}
		if match.AwayTeam.Name == "" {
			t.Error("Match AwayTeam.Name should not be empty")
		}
	}
}

func TestMatchHandler_GetMatch(t *testing.T) {
	router := setupMatchHandlerRouter(t)

	tests := []struct {
		name       string
		matchID    string
		wantStatus int
	}{
		{
			name:       "valid match ID 1",
			matchID:    "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid match ID 2",
			matchID:    "2",
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid match ID 5",
			matchID:    "5",
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existent match ID",
			matchID:    "999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid match ID format",
			matchID:    "invalid",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/betting/matches/"+tt.matchID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var match model.Match
				if err := json.Unmarshal(w.Body.Bytes(), &match); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if match.ID.String() == "" || match.ID.String() == "00000000-0000-0000-0000-000000000000" {
					t.Error("Match ID should be a valid UUID")
				}
				if match.League == "" {
					t.Error("Match League should not be empty")
				}
				if match.HomeTeam.Name == "" {
					t.Error("Match HomeTeam.Name should not be empty")
				}
				if match.AwayTeam.Name == "" {
					t.Error("Match AwayTeam.Name should not be empty")
				}
			}

			if tt.wantStatus == http.StatusNotFound {
				var errResp ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
					t.Fatalf("Failed to unmarshal error response: %v", err)
				}
				if errResp.Error != "match not found" {
					t.Errorf("Expected error message 'match not found', got '%s'", errResp.Error)
				}
			}
		})
	}
}

func TestMatchHandler_GetMatch_Details(t *testing.T) {
	router := setupMatchHandlerRouter(t)

	// Test getting match ID 1 (Manchester United vs Manchester City)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/betting/matches/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var match model.Match
	if err := json.Unmarshal(w.Body.Bytes(), &match); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if match.League != "Premier League" {
		t.Errorf("Expected League 'Premier League', got '%s'", match.League)
	}

	if match.HomeTeam.Name != "Manchester United" {
		t.Errorf("Expected HomeTeam.Name 'Manchester United', got '%s'", match.HomeTeam.Name)
	}

	if match.AwayTeam.Name != "Manchester City" {
		t.Errorf("Expected AwayTeam.Name 'Manchester City', got '%s'", match.AwayTeam.Name)
	}

	if match.Venue != "Old Trafford" {
		t.Errorf("Expected Venue 'Old Trafford', got '%s'", match.Venue)
	}

	if match.Status != "scheduled" {
		t.Errorf("Expected Status 'scheduled', got '%s'", match.Status)
	}

	// Verify team details are populated
	if match.HomeTeam.Country != "England" {
		t.Errorf("Expected HomeTeam.Country 'England', got '%s'", match.HomeTeam.Country)
	}

	if match.AwayTeam.Elo <= 0 {
		t.Error("Expected AwayTeam.Elo to be positive")
	}
}

func TestMatchHandler_GetMatchOdds(t *testing.T) {
	router := setupMatchHandlerRouter(t)

	tests := []struct {
		name          string
		matchID       string
		wantStatus    int
		expectedCount int
	}{
		{
			name:          "match 1 with 5 odds",
			matchID:       "1",
			wantStatus:    http.StatusOK,
			expectedCount: 5,
		},
		{
			name:          "match 2 with 3 odds",
			matchID:       "2",
			wantStatus:    http.StatusOK,
			expectedCount: 3,
		},
		{
			name:          "match 3 with 3 odds",
			matchID:       "3",
			wantStatus:    http.StatusOK,
			expectedCount: 3,
		},
		{
			name:          "match 4 with 3 odds",
			matchID:       "4",
			wantStatus:    http.StatusOK,
			expectedCount: 3,
		},
		{
			name:          "match 5 with 3 odds",
			matchID:       "5",
			wantStatus:    http.StatusOK,
			expectedCount: 3,
		},
		{
			name:          "non-existent match returns empty array",
			matchID:       "999",
			wantStatus:    http.StatusOK,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/betting/matches/"+tt.matchID+"/odds", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			var odds []model.Odds
			if err := json.Unmarshal(w.Body.Bytes(), &odds); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if len(odds) != tt.expectedCount {
				t.Errorf("Expected %d odds, got %d", tt.expectedCount, len(odds))
			}

			// Verify each odd has required fields
			for _, odd := range odds {
				if odd.ID.String() == "" || odd.ID.String() == "00000000-0000-0000-0000-000000000000" {
					t.Error("Odds ID should be a valid UUID")
				}
				if odd.Bookmaker == "" {
					t.Error("Odds Bookmaker should not be empty")
				}
				if odd.Market == "" {
					t.Error("Odds Market should not be empty")
				}
				if odd.Price <= 0 {
					t.Errorf("Odds Price should be positive, got %f", odd.Price)
				}
			}
		})
	}
}

func TestMatchHandler_GetMatchOdds_Details(t *testing.T) {
	router := setupMatchHandlerRouter(t)

	// Get odds for match 1
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/betting/matches/1/odds", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var odds []model.Odds
	if err := json.Unmarshal(w.Body.Bytes(), &odds); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check that we have both 1X2 and O/U 2.5 markets
	markets := make(map[string]bool)
	outcomes := make(map[string]bool)
	for _, odd := range odds {
		markets[odd.Market] = true
		outcomes[odd.Outcome] = true

		// All odds should be from Bet365 based on mock data
		if odd.Bookmaker != "Bet365" {
			t.Errorf("Expected Bookmaker 'Bet365', got '%s'", odd.Bookmaker)
		}
	}

	if !markets["1X2"] {
		t.Error("Expected to find 1X2 market in odds")
	}

	if !markets["O/U 2.5"] {
		t.Error("Expected to find O/U 2.5 market in odds")
	}

	// Check 1X2 outcomes
	if !outcomes["1"] {
		t.Error("Expected to find outcome '1' (home win) in odds")
	}
	if !outcomes["X"] {
		t.Error("Expected to find outcome 'X' (draw) in odds")
	}
	if !outcomes["2"] {
		t.Error("Expected to find outcome '2' (away win) in odds")
	}
}

func TestMatchHandler_RoutesRegistered(t *testing.T) {
	router := setupMatchHandlerRouter(t)

	// Test that all routes are accessible (not 404 for unknown route)
	routes := []string{
		"/api/v1/betting/matches",
		"/api/v1/betting/matches/1",
		"/api/v1/betting/matches/1/odds",
	}

	for _, route := range routes {
		t.Run(route, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, route, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should not get 404 (route not found) - acceptable codes are 200, 400, 500
			// A 404 would indicate the route wasn't registered
			if w.Code == http.StatusNotFound {
				// Check if it's a "match not found" error or route not found
				var errResp ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
					t.Errorf("Route %s seems to not be registered", route)
				} else if errResp.Error != "match not found" {
					t.Errorf("Route %s seems to not be registered", route)
				}
			}
		})
	}
}
