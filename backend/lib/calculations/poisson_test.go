package calculations

import (
	"math"
	"testing"
)

func TestPoissonProbability(t *testing.T) {
	tests := []struct {
		name     string
		lambda   float64
		k        int
		expected float64
		delta    float64
	}{
		{"lambda=2, k=0", 2, 0, 0.1353, 0.001},
		{"lambda=2, k=1", 2, 1, 0.2707, 0.001},
		{"lambda=2, k=2", 2, 2, 0.2707, 0.001},
		{"lambda=2, k=3", 2, 3, 0.1804, 0.001},
		{"lambda=0, k=0", 0, 0, 0, 0.001},      // Edge case
		{"lambda=2, k=-1", 2, -1, 0, 0.001},    // Negative k
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PoissonProbability(tt.lambda, tt.k)
			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("PoissonProbability(%v, %v) = %v, expected %v (±%v)",
					tt.lambda, tt.k, result, tt.expected, tt.delta)
			}
		})
	}
}

func TestCalculatePoissonPrediction(t *testing.T) {
	// Test with typical match stats
	result := CalculatePoissonPrediction(1.8, 1.0, 1.5, 1.2, 2.75)

	// Check that probabilities sum to approximately 100%
	totalProb := result.HomeWinProb + result.DrawProb + result.AwayWinProb
	if math.Abs(totalProb-100) > 1 {
		t.Errorf("Total probability = %v, expected ~100", totalProb)
	}

	// Check that over/under sum to approximately 100%
	overUnderTotal := result.Over25Prob + result.Under25Prob
	if math.Abs(overUnderTotal-100) > 1 {
		t.Errorf("Over/Under total = %v, expected ~100", overUnderTotal)
	}

	// Home team with better stats should have higher win probability
	if result.HomeWinProb < result.AwayWinProb {
		t.Error("Home team with better attack should have higher win probability")
	}

	// Expected goals should be positive
	if result.ExpectedHomeGoals <= 0 || result.ExpectedAwayGoals <= 0 {
		t.Error("Expected goals should be positive")
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		n        int
		expected float64
	}{
		{0, 1},
		{1, 1},
		{5, 120},
		{10, 3628800},
	}

	for _, tt := range tests {
		result := Factorial(tt.n)
		if result != tt.expected {
			t.Errorf("Factorial(%d) = %v, expected %v", tt.n, result, tt.expected)
		}
	}
}

func TestCalculateGoalProbabilities(t *testing.T) {
	homeWin, draw, awayWin, over25, under25, btts := CalculateGoalProbabilities(1.5, 1.2)

	// Check probabilities are valid
	if homeWin < 0 || homeWin > 1 {
		t.Errorf("homeWin = %v, should be 0-1", homeWin)
	}
	if draw < 0 || draw > 1 {
		t.Errorf("draw = %v, should be 0-1", draw)
	}
	if awayWin < 0 || awayWin > 1 {
		t.Errorf("awayWin = %v, should be 0-1", awayWin)
	}

	// Check that 1X2 sums to ~1
	total := homeWin + draw + awayWin
	if math.Abs(total-1) > 0.01 {
		t.Errorf("1X2 total = %v, expected ~1", total)
	}

	// Check over/under
	if math.Abs((over25+under25)-1) > 0.01 {
		t.Errorf("Over/Under total = %v, expected ~1", over25+under25)
	}

	// BTTS should be in valid range
	if btts < 0 || btts > 1 {
		t.Errorf("btts = %v, should be 0-1", btts)
	}
}
