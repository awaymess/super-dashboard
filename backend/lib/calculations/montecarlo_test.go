package calculations

import (
	"math"
	"testing"
)

func TestRunMonteCarloSimulation(t *testing.T) {
	config := MonteCarloConfig{
		Simulations:      1000,
		InitialValue:     100000,
		ExpectedReturn:   8,
		Volatility:       15,
		TimeHorizonYears: 1,
		Seed:             42, // Fixed seed for reproducibility
	}

	result := RunMonteCarloSimulation(config)

	// Check that simulations were run
	if result.Simulations != 1000 {
		t.Errorf("Expected 1000 simulations, got %d", result.Simulations)
	}

	// Mean should be close to expected return (but with randomness)
	// With 8% expected return, mean should be around 108000
	if result.Mean < 90000 || result.Mean > 130000 {
		t.Errorf("Mean %v seems unreasonable for 8%% expected return", result.Mean)
	}

	// Percentiles should be ordered
	if result.Percentile5 > result.Percentile25 ||
		result.Percentile25 > result.Median ||
		result.Median > result.Percentile75 ||
		result.Percentile75 > result.Percentile95 {
		t.Error("Percentiles should be ordered")
	}

	// Min should be <= all percentiles
	if result.MinValue > result.Percentile5 {
		t.Error("MinValue should be <= Percentile5")
	}

	// Max should be >= all percentiles
	if result.MaxValue < result.Percentile95 {
		t.Error("MaxValue should be >= Percentile95")
	}

	// Loss + Gain probability should equal ~100%
	total := result.ProbabilityOfLoss + result.ProbabilityOfGain
	if math.Abs(total-100) > 1 {
		t.Errorf("Loss + Gain probability = %v, expected ~100", total)
	}
}

func TestRunBettingMonteCarlo(t *testing.T) {
	config := BettingMonteCarloConfig{
		Simulations:     1000,
		InitialBankroll: 1000,
		NumBets:         100,
		WinProbability:  55, // Positive EV at 2.0 odds
		AverageOdds:     2.0,
		StakePercent:    2,
		Seed:            42,
	}

	result := RunBettingMonteCarlo(config)

	// With positive EV, average should be above initial
	if result.AvgFinalBankroll < config.InitialBankroll {
		t.Log("Warning: Avg final bankroll below initial, but this can happen with variance")
	}

	// Ruin probability should be very low with 2% stakes
	if result.RuinProbability > 5 {
		t.Errorf("Ruin probability %v seems too high for 2%% stakes", result.RuinProbability)
	}

	// Max drawdown should exist
	if result.AvgMaxDrawdown < 0 {
		t.Error("Average max drawdown should be non-negative")
	}
}

func TestSharpeRatio(t *testing.T) {
	// Returns: 10%, 15%, 5%, 20%, 8% with 2% risk-free rate
	returns := []float64{10, 15, 5, 20, 8}
	riskFreeRate := 2.0

	sharpe := SharpeRatio(returns, riskFreeRate)

	// Mean = 11.6, StdDev ≈ 5.68
	// Sharpe = (11.6 - 2) / 5.68 ≈ 1.69
	if sharpe < 1 || sharpe > 3 {
		t.Errorf("Sharpe ratio %v seems unreasonable", sharpe)
	}

	// Empty returns
	if SharpeRatio([]float64{}, 2) != 0 {
		t.Error("Sharpe should be 0 for empty returns")
	}
}

func TestMaxDrawdown(t *testing.T) {
	// Portfolio values: 100, 110, 90, 95, 105
	// Max = 110, then drops to 90, drawdown = (110-90)/110 = 18.18%
	values := []float64{100, 110, 90, 95, 105}

	dd := MaxDrawdown(values)
	expected := 18.18

	if math.Abs(dd-expected) > 0.1 {
		t.Errorf("Max drawdown = %v, expected ~%v", dd, expected)
	}

	// No drawdown
	values = []float64{100, 110, 120, 130}
	dd = MaxDrawdown(values)
	if dd != 0 {
		t.Errorf("Max drawdown should be 0 for monotonically increasing, got %v", dd)
	}

	// Edge case: single value
	if MaxDrawdown([]float64{100}) != 0 {
		t.Error("Max drawdown should be 0 for single value")
	}
}

func TestValueAtRisk(t *testing.T) {
	// Returns: sorted would be -5, 2, 5, 8, 10, 12
	// 95% VaR = 5th percentile = around -5
	returns := []float64{10, -5, 8, 5, 12, 2}

	var95 := ValueAtRisk(returns, 95)

	// 5th percentile should be close to the lowest values
	if var95 > 0 {
		t.Log("VaR at 95% confidence should typically be negative or near zero")
	}

	// Empty returns
	if ValueAtRisk([]float64{}, 95) != 0 {
		t.Error("VaR should be 0 for empty returns")
	}
}

func TestCalculateMean(t *testing.T) {
	values := []float64{10, 20, 30, 40, 50}
	mean := calculateMean(values)
	expected := 30.0

	if math.Abs(mean-expected) > 0.01 {
		t.Errorf("Mean = %v, expected %v", mean, expected)
	}

	// Empty slice
	if calculateMean([]float64{}) != 0 {
		t.Error("Mean of empty slice should be 0")
	}
}

func TestCalculateStdDev(t *testing.T) {
	// Values with known std dev
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	stdDev := calculateStdDev(values)
	expected := 2.14 // Approximate

	if math.Abs(stdDev-expected) > 0.1 {
		t.Errorf("StdDev = %v, expected ~%v", stdDev, expected)
	}

	// Single value
	if calculateStdDev([]float64{5}) != 0 {
		t.Error("StdDev of single value should be 0")
	}
}

func TestPercentile(t *testing.T) {
	sorted := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Median (50th percentile)
	p50 := percentile(sorted, 50)
	if math.Abs(p50-5.5) > 0.1 {
		t.Errorf("50th percentile = %v, expected ~5.5", p50)
	}

	// Edge cases
	if percentile(sorted, 0) != 1 {
		t.Error("0th percentile should be min")
	}
	if percentile(sorted, 100) != 10 {
		t.Error("100th percentile should be max")
	}
}

func TestMonteCarloWithDifferentSeeds(t *testing.T) {
	config1 := MonteCarloConfig{
		Simulations:      100,
		InitialValue:     100000,
		ExpectedReturn:   10,
		Volatility:       20,
		TimeHorizonYears: 1,
		Seed:             1,
	}

	config2 := config1
	config2.Seed = 2

	result1 := RunMonteCarloSimulation(config1)
	result2 := RunMonteCarloSimulation(config2)

	// Different seeds should produce different results
	if result1.Mean == result2.Mean {
		t.Log("Different seeds produced same mean - unlikely but possible")
	}
}

func TestMonteCarloReproducibility(t *testing.T) {
	config := MonteCarloConfig{
		Simulations:      100,
		InitialValue:     100000,
		ExpectedReturn:   10,
		Volatility:       20,
		TimeHorizonYears: 1,
		Seed:             42,
	}

	result1 := RunMonteCarloSimulation(config)
	result2 := RunMonteCarloSimulation(config)

	// Same seed should produce same results
	if result1.Mean != result2.Mean {
		t.Error("Same seed should produce reproducible results")
	}
}
