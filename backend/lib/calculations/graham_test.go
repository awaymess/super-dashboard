package calculations

import (
	"math"
	"testing"
)

func TestCalculateGrahamNumber(t *testing.T) {
	// Classic Graham Number: sqrt(22.5 * EPS * Book Value)
	// EPS = 5, BV = 40 -> sqrt(22.5 * 5 * 40) = sqrt(4500) = 67.08
	result := CalculateGrahamNumber(5, 40)
	expected := 67.08

	if math.Abs(result-expected) > 0.1 {
		t.Errorf("Graham Number = %v, expected %v", result, expected)
	}

	// Edge cases
	if CalculateGrahamNumber(-5, 40) != 0 {
		t.Error("Should return 0 for negative EPS")
	}
	if CalculateGrahamNumber(5, -40) != 0 {
		t.Error("Should return 0 for negative book value")
	}
}

func TestCalculateModifiedGrahamValue(t *testing.T) {
	// V = EPS × (8.5 + 2g) × 4.4 / Y
	// EPS = 5, g = 10%, Y = 4.4%
	// = 5 * (8.5 + 20) * 4.4 / 4.4 = 5 * 28.5 = 142.5
	result := CalculateModifiedGrahamValue(5, 10, 4.4)
	expected := 142.5

	if math.Abs(result-expected) > 0.1 {
		t.Errorf("Modified Graham Value = %v, expected %v", result, expected)
	}

	// No growth
	result = CalculateModifiedGrahamValue(5, 0, 4.4)
	expected = 42.5 // 5 * 8.5 = 42.5
	if math.Abs(result-expected) > 0.1 {
		t.Errorf("Modified Graham Value (no growth) = %v, expected %v", result, expected)
	}

	// Edge case: negative EPS
	if CalculateModifiedGrahamValue(-5, 10, 4.4) != 0 {
		t.Error("Should return 0 for negative EPS")
	}
}

func TestCalculateGrahamAnalysis(t *testing.T) {
	inputs := GrahamInputs{
		EPS:               5,
		BookValuePerShare: 40,
		CurrentPrice:      50,
		GrowthRate:        10,
		AAAYield:          4.4,
	}

	result := CalculateGrahamAnalysis(inputs)

	// Check PE ratio
	expectedPE := 50.0 / 5.0
	if math.Abs(result.Analysis.PERatio-expectedPE) > 0.01 {
		t.Errorf("PE = %v, expected %v", result.Analysis.PERatio, expectedPE)
	}

	// Check PB ratio
	expectedPB := 50.0 / 40.0
	if math.Abs(result.Analysis.PBRatio-expectedPB) > 0.01 {
		t.Errorf("PB = %v, expected %v", result.Analysis.PBRatio, expectedPB)
	}

	// PEG should be calculated when growth > 0
	if result.Analysis.PEG == nil {
		t.Error("PEG should be calculated when growth rate > 0")
	} else {
		expectedPEG := expectedPE / 10.0
		if math.Abs(*result.Analysis.PEG-expectedPEG) > 0.01 {
			t.Errorf("PEG = %v, expected %v", *result.Analysis.PEG, expectedPEG)
		}
	}

	// Check defensive criteria
	// PE = 10 <= 15, PB = 1.25 <= 1.5, PE*PB = 12.5 <= 22.5
	if !result.Analysis.IsDefensive {
		t.Error("Stock should meet defensive criteria")
	}

	// Margin of safety should be positive (intrinsic > price)
	if result.MarginOfSafety <= 0 {
		t.Error("Margin of safety should be positive for undervalued stock")
	}
}

func TestScreenDefensiveStocks(t *testing.T) {
	stocks := []StockForScreening{
		{Symbol: "GOOD", EPS: 5, BookValue: 40, Price: 50},  // PE=10, PB=1.25, PE*PB=12.5 ✓
		{Symbol: "HIGHPE", EPS: 5, BookValue: 40, Price: 100}, // PE=20, PB=2.5 ✗
		{Symbol: "LOWEPS", EPS: 0, BookValue: 40, Price: 50},  // Zero EPS ✗
		{Symbol: "OK", EPS: 10, BookValue: 100, Price: 100},   // PE=10, PB=1, PE*PB=10 ✓
	}

	result := ScreenDefensiveStocks(stocks)

	if len(result) != 2 {
		t.Errorf("Expected 2 defensive stocks, got %d", len(result))
	}

	// Check that GOOD and OK are in the result
	found := map[string]bool{}
	for _, s := range result {
		found[s] = true
	}
	if !found["GOOD"] {
		t.Error("GOOD should be in defensive stocks")
	}
	if !found["OK"] {
		t.Error("OK should be in defensive stocks")
	}
}

func TestCalculateNCAV(t *testing.T) {
	// Current Assets = 100M, Total Liabilities = 60M, Shares = 10M
	result := CalculateNCAV(100000000, 60000000, 10000000)

	expectedNCAV := 40000000.0
	expectedPerShare := 4.0

	if math.Abs(result.NCAV-expectedNCAV) > 0.01 {
		t.Errorf("NCAV = %v, expected %v", result.NCAV, expectedNCAV)
	}
	if math.Abs(result.NCAVPerShare-expectedPerShare) > 0.01 {
		t.Errorf("NCAV per share = %v, expected %v", result.NCAVPerShare, expectedPerShare)
	}

	// Edge case: zero shares
	result = CalculateNCAV(100000000, 60000000, 0)
	if result.NCAVPerShare != 0 {
		t.Error("NCAV per share should be 0 when shares outstanding is 0")
	}
}

func TestIsNetNet(t *testing.T) {
	// Net-net: Market cap < 2/3 of NCAV
	// NCAV = 100M - 60M = 40M
	// 2/3 of 40M = 26.67M

	// Market cap 20M < 26.67M -> Net-net
	if !IsNetNet(100000000, 60000000, 20000000) {
		t.Error("Should be a net-net bargain")
	}

	// Market cap 30M > 26.67M -> Not net-net
	if IsNetNet(100000000, 60000000, 30000000) {
		t.Error("Should not be a net-net bargain")
	}
}

func TestBuffettIntrinsicValue(t *testing.T) {
	// Owner earnings 10M, growth 10%, discount 12%, 10 years
	result := BuffettIntrinsicValue(10000000, 10, 12, 10)

	// Should return positive value
	if result <= 0 {
		t.Error("Intrinsic value should be positive")
	}

	// Invalid inputs (growth >= discount)
	if BuffettIntrinsicValue(10000000, 15, 10, 10) != 0 {
		t.Error("Should return 0 when growth >= discount")
	}
}

func TestPBVValuation(t *testing.T) {
	// Book value $50, fair P/BV 1.5 -> $75
	result := PBVValuation(50, 1.5)
	if math.Abs(result-75) > 0.01 {
		t.Errorf("PBV Valuation = %v, expected 75", result)
	}
}

func TestPEValuation(t *testing.T) {
	// EPS $5, fair P/E 15 -> $75
	result := PEValuation(5, 15)
	if math.Abs(result-75) > 0.01 {
		t.Errorf("PE Valuation = %v, expected 75", result)
	}
}

func TestGrahamRating(t *testing.T) {
	// High margin of safety -> undervalued
	inputs := GrahamInputs{
		EPS:               10,
		BookValuePerShare: 80,
		CurrentPrice:      50,  // Way below intrinsic
		GrowthRate:        5,
		AAAYield:          4.4,
	}
	result := CalculateGrahamAnalysis(inputs)
	if result.Rating != "undervalued" {
		t.Errorf("Expected 'undervalued', got '%s' (MoS: %v)", result.Rating, result.MarginOfSafety)
	}

	// Stock priced way above intrinsic value -> overvalued
	// Graham Number = sqrt(22.5 * 10 * 80) = sqrt(18000) = ~134.16
	// Modified Graham Value = 10 * (8.5 + 10) * 4.4 / 4.4 = 185
	// Intrinsic = max(134.16, 185) = 185
	// At price 500: MoS = (185-500)/185 = -170% -> definitely overvalued
	inputs.CurrentPrice = 500
	result = CalculateGrahamAnalysis(inputs)
	if result.Rating != "overvalued" {
		t.Errorf("Expected 'overvalued', got '%s' (MoS: %v)", result.Rating, result.MarginOfSafety)
	}
}
