package calculations

import (
	"math"
	"testing"
)

func TestCalculateExpectedScore(t *testing.T) {
	// Equal ratings should give 50%
	result := CalculateExpectedScore(1500, 1500)
	if math.Abs(result-0.5) > 0.001 {
		t.Errorf("Equal ratings should give 0.5, got %v", result)
	}

	// Higher rated should have higher expected score
	result = CalculateExpectedScore(1600, 1500)
	if result <= 0.5 {
		t.Errorf("Higher rated should have >50%% expected score, got %v", result)
	}

	// 400 point difference should give ~91% expected score
	result = CalculateExpectedScore(1900, 1500)
	expected := 0.909
	if math.Abs(result-expected) > 0.01 {
		t.Errorf("400 point difference should give ~91%%, got %v", result)
	}
}

func TestCalculateNewRating(t *testing.T) {
	// Win when expected to win slightly
	newRating := CalculateNewRating(1500, 0.6, 1.0, 32)
	// Should increase: 1500 + 32*(1-0.6) = 1500 + 12.8
	expected := 1512.8
	if math.Abs(newRating-expected) > 0.1 {
		t.Errorf("Expected %v, got %v", expected, newRating)
	}

	// Lose when expected to win
	newRating = CalculateNewRating(1500, 0.6, 0.0, 32)
	// Should decrease: 1500 + 32*(0-0.6) = 1500 - 19.2
	expected = 1480.8
	if math.Abs(newRating-expected) > 0.1 {
		t.Errorf("Expected %v, got %v", expected, newRating)
	}
}

func TestCalculateELOMatchProbabilities(t *testing.T) {
	probs := CalculateELOMatchProbabilities(1600, 1400, 100)

	// Home team should be favored
	if probs.HomeWin <= probs.AwayWin {
		t.Error("Higher rated home team should be favored")
	}

	// Probabilities should sum to ~100%
	total := probs.HomeWin + probs.Draw + probs.AwayWin
	if math.Abs(total-100) > 0.1 {
		t.Errorf("Total probability = %v, expected ~100", total)
	}

	// All probabilities should be positive
	if probs.HomeWin < 0 || probs.Draw < 0 || probs.AwayWin < 0 {
		t.Error("All probabilities should be positive")
	}
}

func TestUpdateRatings(t *testing.T) {
	// Home win
	result := UpdateRatings(1500, 1500, 2, 1, 32, 100)
	
	if result.HomeRating.Change <= 0 {
		t.Error("Winner should gain rating")
	}
	if result.AwayRating.Change >= 0 {
		t.Error("Loser should lose rating")
	}

	// Draw
	result = UpdateRatings(1500, 1500, 1, 1, 32, 100)
	// With home advantage, home team is favored, so draw means home loses points
	// and away gains points
	if result.HomeRating.Change >= 0 {
		t.Error("Home team with advantage drawing should lose rating")
	}

	// Big win should cause bigger rating change
	result1 := UpdateRatings(1500, 1500, 1, 0, 32, 100)
	result3 := UpdateRatings(1500, 1500, 4, 0, 32, 100)
	if result1.HomeRating.Change >= result3.HomeRating.Change {
		t.Error("Bigger margin should cause bigger rating change")
	}
}

func TestSimulateELOSeason(t *testing.T) {
	teams := map[string]float64{
		"TeamA": 1500,
		"TeamB": 1500,
		"TeamC": 1500,
	}

	matches := []MatchResult{
		{HomeTeam: "TeamA", AwayTeam: "TeamB", HomeScore: 2, AwayScore: 1},
		{HomeTeam: "TeamB", AwayTeam: "TeamC", HomeScore: 1, AwayScore: 1},
		{HomeTeam: "TeamC", AwayTeam: "TeamA", HomeScore: 0, AwayScore: 3},
	}

	ratings := SimulateELOSeason(teams, matches, 32)

	// TeamA won both, should have highest rating
	if ratings["TeamA"] <= ratings["TeamB"] || ratings["TeamA"] <= ratings["TeamC"] {
		t.Error("TeamA should have highest rating after winning both matches")
	}

	// Total ELO should be conserved (approximately)
	totalBefore := 1500.0 * 3
	totalAfter := ratings["TeamA"] + ratings["TeamB"] + ratings["TeamC"]
	if math.Abs(totalBefore-totalAfter) > 1 {
		t.Errorf("Total ELO not conserved: before=%v, after=%v", totalBefore, totalAfter)
	}
}

func TestRatingToTier(t *testing.T) {
	tests := []struct {
		rating   int
		expected string
	}{
		{2100, "Elite"},
		{1900, "Strong"},
		{1700, "Above Average"},
		{1500, "Average"},
		{1300, "Below Average"},
		{1100, "Weak"},
	}

	for _, tt := range tests {
		result := RatingToTier(tt.rating)
		if result != tt.expected {
			t.Errorf("RatingToTier(%d) = %s, expected %s", tt.rating, result, tt.expected)
		}
	}
}

func TestGetInitialRating(t *testing.T) {
	if GetInitialRating() != 1500 {
		t.Errorf("Initial rating should be 1500")
	}
}

func TestPredictMatchOutcome(t *testing.T) {
	// Team with good form should have better odds
	goodForm := PredictMatchOutcome(1500, 1500, 0.8, 0.3)
	badForm := PredictMatchOutcome(1500, 1500, 0.3, 0.8)

	if goodForm.HomeWin <= badForm.HomeWin {
		t.Error("Better home form should improve home win probability")
	}
}
