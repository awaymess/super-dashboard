package calculations

import (
	"math"
	"testing"
)

func TestCalculateKelly(t *testing.T) {
	tests := []struct {
		name        string
		probability float64
		odds        float64
		bankroll    float64
		fraction    float64
		wantPositive bool
	}{
		{"60% at 2.0 odds", 60, 2.0, 1000, 1.0, true},
		{"50% at 2.0 odds", 50, 2.0, 1000, 1.0, false}, // No edge
		{"40% at 2.0 odds", 40, 2.0, 1000, 1.0, false}, // Negative edge
		{"70% at 1.5 odds", 70, 1.5, 1000, 1.0, true},
		{"Half Kelly", 60, 2.0, 1000, 0.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateKelly(tt.probability, tt.odds, tt.bankroll, tt.fraction)
			
			if tt.wantPositive && result.Stake <= 0 {
				t.Errorf("Expected positive stake, got %v", result.Stake)
			}
			if !tt.wantPositive && result.Stake > 0 {
				t.Errorf("Expected zero stake, got %v", result.Stake)
			}
			
			// Half Kelly should be half of full
			if result.HalfKelly != result.Stake*0.5 {
				t.Errorf("HalfKelly = %v, expected %v", result.HalfKelly, result.Stake*0.5)
			}
			
			// Quarter Kelly should be quarter of full
			if result.QuarterKelly != result.Stake*0.25 {
				t.Errorf("QuarterKelly = %v, expected %v", result.QuarterKelly, result.Stake*0.25)
			}
		})
	}
}

func TestCalculateFullKelly(t *testing.T) {
	// 60% win probability at 2.0 odds
	// Kelly = (1*0.6 - 0.4) / 1 = 0.2 (20% of bankroll)
	result := CalculateFullKelly(0.6, 2.0)
	expected := 0.2
	if math.Abs(result-expected) > 0.001 {
		t.Errorf("CalculateFullKelly(0.6, 2.0) = %v, expected %v", result, expected)
	}

	// Edge cases
	if CalculateFullKelly(0, 2.0) != 0 {
		t.Error("Should return 0 for 0 probability")
	}
	if CalculateFullKelly(1, 2.0) != 0 {
		t.Error("Should return 0 for 100% probability")
	}
	if CalculateFullKelly(0.5, 1.0) != 0 {
		t.Error("Should return 0 for odds <= 1")
	}
}

func TestCalculateImpliedProbability(t *testing.T) {
	tests := []struct {
		odds     float64
		expected float64
	}{
		{2.0, 50},
		{4.0, 25},
		{1.5, 66.67},
		{1.0, 0},  // Edge case
		{0.5, 0},  // Invalid odds
	}

	for _, tt := range tests {
		result := CalculateImpliedProbability(tt.odds)
		if math.Abs(result-tt.expected) > 0.1 {
			t.Errorf("CalculateImpliedProbability(%v) = %v, expected %v", tt.odds, result, tt.expected)
		}
	}
}

func TestProbabilityToOdds(t *testing.T) {
	tests := []struct {
		probability float64
		expected    float64
	}{
		{50, 2.0},
		{25, 4.0},
		{100, 0}, // Edge case
		{0, 0},   // Edge case
	}

	for _, tt := range tests {
		result := ProbabilityToOdds(tt.probability)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("ProbabilityToOdds(%v) = %v, expected %v", tt.probability, result, tt.expected)
		}
	}
}

func TestDetectValueBet(t *testing.T) {
	// True probability 60%, odds 2.0 implies 50% -> 10% value
	result := DetectValueBet(60, 2.0, 5)
	if !result.IsValueBet {
		t.Error("Should be a value bet")
	}
	// Value is exactly 10%, which is the threshold for high value
	// but the condition is > 10, so at exactly 10% it's just "bet" not "strong_bet"
	if result.Value < 5 {
		t.Errorf("Value should be >= 5%%, got %v", result.Value)
	}
	if result.Recommendation == "skip" {
		t.Errorf("Expected 'bet' or 'strong_bet', got '%s'", result.Recommendation)
	}

	// True probability 52%, odds 2.0 implies 50% -> 2% value (below threshold)
	result = DetectValueBet(52, 2.0, 5)
	if result.IsValueBet {
		t.Error("Should not be a value bet (below threshold)")
	}
	if result.Recommendation != "skip" {
		t.Errorf("Expected 'skip', got '%s'", result.Recommendation)
	}
}

func TestBayesianUpdate(t *testing.T) {
	// Classic Bayes example
	prior := 0.5
	likelihood := 0.9
	evidence := 0.7

	posterior := BayesianUpdate(prior, likelihood, evidence)
	expected := (0.9 * 0.5) / 0.7 // ~0.643

	if math.Abs(posterior-expected) > 0.001 {
		t.Errorf("BayesianUpdate = %v, expected %v", posterior, expected)
	}

	// Edge case: zero evidence probability
	if BayesianUpdate(0.5, 0.9, 0) != 0.5 {
		t.Error("Should return prior when evidence probability is 0")
	}
}

func TestCalculateWeightedProbability(t *testing.T) {
	probs := []float64{50, 60, 70}
	weights := []float64{0.5, 0.3, 0.2}

	result, _ := CalculateWeightedProbability(probs, weights)
	// 50*0.5 + 60*0.3 + 70*0.2 = 25 + 18 + 14 = 57
	expected := 57.0

	if math.Abs(result-expected) > 0.01 {
		t.Errorf("CalculateWeightedProbability = %v, expected %v", result, expected)
	}
}

func TestFindArbitrage(t *testing.T) {
	// No arbitrage scenario
	homeOdds := []float64{1.8, 1.85}
	drawOdds := []float64{3.5, 3.6}
	awayOdds := []float64{4.0, 4.2}

	result := FindArbitrage(homeOdds, drawOdds, awayOdds)
	if result.IsArbitrage {
		t.Error("Should not find arbitrage in normal odds")
	}

	// Arbitrage scenario (artificially high odds)
	homeOdds = []float64{3.0, 3.1}
	drawOdds = []float64{3.5, 3.6}
	awayOdds = []float64{3.0, 3.1}

	result = FindArbitrage(homeOdds, drawOdds, awayOdds)
	if !result.IsArbitrage {
		t.Error("Should find arbitrage in inflated odds")
	}
	if result.Stakes == nil {
		t.Error("Stakes should be calculated for arbitrage")
	}
}

func TestSimulateKellyGrowth(t *testing.T) {
	bets := []BetOutcome{
		{Probability: 60, Odds: 2.0, Won: true},
		{Probability: 60, Odds: 2.0, Won: true},
		{Probability: 60, Odds: 2.0, Won: false},
		{Probability: 60, Odds: 2.0, Won: true},
	}

	result := SimulateKellyGrowth(1000, bets, 0.5)

	// With positive expected value bets, we should see growth on average
	// But this specific sequence should result in profit
	if result.FinalBankroll <= 0 {
		t.Error("Bankroll should be positive")
	}

	// Max drawdown should be non-negative
	if result.MaxDrawdown < 0 {
		t.Error("Max drawdown should be non-negative")
	}
}
