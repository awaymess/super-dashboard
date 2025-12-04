package repository

import (
	"os"
	"path/filepath"
	"testing"
)

// findMockDataPath finds the mock data directory for tests
func findMockDataPath() string {
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

func TestNewMockMatchRepository(t *testing.T) {
	mockPath := findMockDataPath()

	repo, err := NewMockMatchRepository(mockPath)
	if err != nil {
		t.Fatalf("Failed to create mock match repository: %v", err)
	}

	if repo == nil {
		t.Fatal("Expected repository to be non-nil")
	}
}

func TestNewMockMatchRepository_FileNotFound(t *testing.T) {
	_, err := NewMockMatchRepository("/nonexistent/path/matches.json")
	if err == nil {
		t.Fatal("Expected error for non-existent file")
	}
}

func TestMockMatchRepository_GetAll(t *testing.T) {
	mockPath := findMockDataPath()

	repo, err := NewMockMatchRepository(mockPath)
	if err != nil {
		t.Fatalf("Failed to create mock match repository: %v", err)
	}

	matches, err := repo.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all matches: %v", err)
	}

	// Based on the mock data, we expect 5 matches
	if len(matches) != 5 {
		t.Errorf("Expected 5 matches, got %d", len(matches))
	}

	// Verify that each match has required fields populated
	for _, match := range matches {
		if match.ID.String() == "" {
			t.Error("Match ID should not be empty")
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
		if match.Venue == "" {
			t.Error("Match Venue should not be empty")
		}
	}
}

func TestMockMatchRepository_GetByID(t *testing.T) {
	mockPath := findMockDataPath()

	repo, err := NewMockMatchRepository(mockPath)
	if err != nil {
		t.Fatalf("Failed to create mock match repository: %v", err)
	}

	tests := []struct {
		name      string
		id        string
		wantErr   bool
		wantError error
	}{
		{
			name:    "valid match ID 1",
			id:      "1",
			wantErr: false,
		},
		{
			name:    "valid match ID 2",
			id:      "2",
			wantErr: false,
		},
		{
			name:    "valid match ID 5",
			id:      "5",
			wantErr: false,
		},
		{
			name:      "non-existent match ID",
			id:        "999",
			wantErr:   true,
			wantError: ErrNotFound,
		},
		{
			name:      "empty match ID",
			id:        "",
			wantErr:   true,
			wantError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := repo.GetByID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.wantError != nil && err != tt.wantError {
					t.Errorf("Expected error %v, got %v", tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if match == nil {
				t.Error("Expected match to be non-nil")
				return
			}

			if match.ID.String() == "" {
				t.Error("Match ID should not be empty")
			}
		})
	}
}

func TestMockMatchRepository_GetByID_MatchDetails(t *testing.T) {
	mockPath := findMockDataPath()

	repo, err := NewMockMatchRepository(mockPath)
	if err != nil {
		t.Fatalf("Failed to create mock match repository: %v", err)
	}

	// Get match ID 1 (Manchester United vs Manchester City)
	match, err := repo.GetByID("1")
	if err != nil {
		t.Fatalf("Failed to get match by ID: %v", err)
	}

	if match.League != "Premier League" {
		t.Errorf("Expected League 'Premier League', got %s", match.League)
	}

	if match.HomeTeam.Name != "Manchester United" {
		t.Errorf("Expected HomeTeam.Name 'Manchester United', got %s", match.HomeTeam.Name)
	}

	if match.AwayTeam.Name != "Manchester City" {
		t.Errorf("Expected AwayTeam.Name 'Manchester City', got %s", match.AwayTeam.Name)
	}

	if match.Venue != "Old Trafford" {
		t.Errorf("Expected Venue 'Old Trafford', got %s", match.Venue)
	}

	if match.Status != "scheduled" {
		t.Errorf("Expected Status 'scheduled', got %s", match.Status)
	}
}

func TestMockMatchRepository_GetOddsByMatchID(t *testing.T) {
	mockPath := findMockDataPath()

	repo, err := NewMockMatchRepository(mockPath)
	if err != nil {
		t.Fatalf("Failed to create mock match repository: %v", err)
	}

	tests := []struct {
		name         string
		matchID      string
		expectedLen  int
		expectOdds   bool
	}{
		{
			name:        "match 1 with 5 odds",
			matchID:     "1",
			expectedLen: 5,
			expectOdds:  true,
		},
		{
			name:        "match 2 with 3 odds",
			matchID:     "2",
			expectedLen: 3,
			expectOdds:  true,
		},
		{
			name:        "match 3 with 3 odds",
			matchID:     "3",
			expectedLen: 3,
			expectOdds:  true,
		},
		{
			name:        "non-existent match returns empty",
			matchID:     "999",
			expectedLen: 0,
			expectOdds:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			odds, err := repo.GetOddsByMatchID(tt.matchID)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(odds) != tt.expectedLen {
				t.Errorf("Expected %d odds, got %d", tt.expectedLen, len(odds))
			}

			if tt.expectOdds && len(odds) > 0 {
				// Verify odds have required fields
				for _, odd := range odds {
					if odd.ID.String() == "" {
						t.Error("Odds ID should not be empty")
					}
					if odd.Bookmaker == "" {
						t.Error("Odds Bookmaker should not be empty")
					}
					if odd.Market == "" {
						t.Error("Odds Market should not be empty")
					}
					if odd.Price <= 0 {
						t.Error("Odds Price should be positive")
					}
				}
			}
		})
	}
}

func TestMockMatchRepository_GetOddsByMatchID_Details(t *testing.T) {
	mockPath := findMockDataPath()

	repo, err := NewMockMatchRepository(mockPath)
	if err != nil {
		t.Fatalf("Failed to create mock match repository: %v", err)
	}

	odds, err := repo.GetOddsByMatchID("1")
	if err != nil {
		t.Fatalf("Failed to get odds: %v", err)
	}

	// Check that we have both 1X2 and O/U 2.5 markets
	markets := make(map[string]bool)
	for _, odd := range odds {
		markets[odd.Market] = true

		// All odds should be from Bet365 based on mock data
		if odd.Bookmaker != "Bet365" {
			t.Errorf("Expected Bookmaker 'Bet365', got %s", odd.Bookmaker)
		}
	}

	if !markets["1X2"] {
		t.Error("Expected to find 1X2 market in odds")
	}

	if !markets["O/U 2.5"] {
		t.Error("Expected to find O/U 2.5 market in odds")
	}
}

func TestStringToUUID(t *testing.T) {
	// Test that the same input always produces the same UUID
	uuid1 := stringToUUID("test-id")
	uuid2 := stringToUUID("test-id")

	if uuid1 != uuid2 {
		t.Errorf("Expected same UUIDs for same input, got %s and %s", uuid1, uuid2)
	}

	// Test that different inputs produce different UUIDs
	uuid3 := stringToUUID("different-id")
	if uuid1 == uuid3 {
		t.Errorf("Expected different UUIDs for different inputs, got same: %s", uuid1)
	}
}
