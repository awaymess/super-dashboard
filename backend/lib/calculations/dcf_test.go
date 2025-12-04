package calculations

import (
	"math"
	"testing"
)

func TestCalculateDCF(t *testing.T) {
	inputs := DCFInputs{
		FreeCashFlow:       100000000, // $100M
		GrowthRate:         10,        // 10%
		TerminalGrowthRate: 3,         // 3%
		DiscountRate:       10,        // 10%
		Years:              5,
		SharesOutstanding:  10000000, // 10M shares
	}

	result := CalculateDCF(inputs)

	// Intrinsic value should be positive
	if result.IntrinsicValue <= 0 {
		t.Error("Intrinsic value should be positive")
	}

	// Per share value should be positive
	if result.PerShareValue <= 0 {
		t.Error("Per share value should be positive")
	}

	// Per share value should equal intrinsic / shares
	expected := result.IntrinsicValue / float64(inputs.SharesOutstanding)
	if math.Abs(result.PerShareValue-expected) > 0.01 {
		t.Errorf("Per share = %v, expected %v", result.PerShareValue, expected)
	}

	// Should have 5 projected cash flows
	if len(result.ProjectedCashFlows) != 5 {
		t.Errorf("Expected 5 projected cash flows, got %d", len(result.ProjectedCashFlows))
	}

	// Cash flows should be increasing with positive growth
	for i := 1; i < len(result.ProjectedCashFlows); i++ {
		if result.ProjectedCashFlows[i].CashFlow <= result.ProjectedCashFlows[i-1].CashFlow {
			t.Error("Cash flows should be increasing with positive growth")
		}
	}

	// Terminal value should be positive
	if result.TerminalValue <= 0 {
		t.Error("Terminal value should be positive")
	}

	// PV of terminal should be less than terminal value (due to discounting)
	if result.PresentValueOfTerminal >= result.TerminalValue {
		t.Error("PV of terminal should be less than terminal value")
	}
}

func TestCalculateWACC(t *testing.T) {
	// 60% equity at 12% cost, 40% debt at 6% cost, 25% tax rate
	wacc := CalculateWACC(60, 40, 12, 6, 25)
	// = 0.6*0.12 + 0.4*0.06*0.75 = 0.072 + 0.018 = 0.09 = 9%
	expected := 0.09

	if math.Abs(wacc-expected) > 0.001 {
		t.Errorf("WACC = %v, expected %v", wacc, expected)
	}
}

func TestCalculateCostOfEquity(t *testing.T) {
	// Risk-free 3%, beta 1.2, market return 10%
	coe := CalculateCostOfEquity(3, 1.2, 10)
	// = 3 + 1.2*(10-3) = 3 + 8.4 = 11.4%
	expected := 11.4

	if math.Abs(coe-expected) > 0.01 {
		t.Errorf("Cost of equity = %v, expected %v", coe, expected)
	}
}

func TestEstimateGrowthRate(t *testing.T) {
	// ROE 15%, retention ratio 60%
	growth := EstimateGrowthRate(15, 60)
	// = 0.15 * 0.6 = 0.09 = 9%
	expected := 9.0

	if math.Abs(growth-expected) > 0.01 {
		t.Errorf("Growth rate = %v, expected %v", growth, expected)
	}
}

func TestReverseDCF(t *testing.T) {
	// Set up a scenario and calculate intrinsic value
	inputs := DCFInputs{
		FreeCashFlow:       100000000,
		GrowthRate:         10,
		TerminalGrowthRate: 3,
		DiscountRate:       10,
		Years:              5,
		SharesOutstanding:  10000000,
	}
	result := CalculateDCF(inputs)
	
	// Now reverse engineer the growth rate
	impliedGrowth := ReverseDCF(
		result.PerShareValue,
		float64(inputs.SharesOutstanding),
		inputs.FreeCashFlow,
		inputs.DiscountRate,
		inputs.TerminalGrowthRate,
		inputs.Years,
	)

	// Should be close to the original growth rate
	if math.Abs(impliedGrowth-inputs.GrowthRate) > 0.5 {
		t.Errorf("Implied growth = %v, expected ~%v", impliedGrowth, inputs.GrowthRate)
	}
}

func TestCalculateMarginOfSafety(t *testing.T) {
	// Stock worth $100, trading at $70 = 30% MoS
	mos := CalculateMarginOfSafety(100, 70)
	if math.Abs(mos-30) > 0.1 {
		t.Errorf("MoS = %v, expected 30", mos)
	}

	// Stock worth $100, trading at $120 = -20% MoS
	mos = CalculateMarginOfSafety(100, 120)
	if math.Abs(mos-(-20)) > 0.1 {
		t.Errorf("MoS = %v, expected -20", mos)
	}

	// Edge case: zero intrinsic value
	mos = CalculateMarginOfSafety(0, 50)
	if mos != 0 {
		t.Error("MoS should be 0 when intrinsic value is 0")
	}
}

func TestGetValuationRating(t *testing.T) {
	tests := []struct {
		mos      float64
		expected ValuationRating
	}{
		{40, Undervalued},
		{30, Undervalued},
		{15, FairValue},
		{0, FairValue},
		{-10, FairValue},
		{-15, Overvalued},
		{-50, Overvalued},
	}

	for _, tt := range tests {
		result := GetValuationRating(tt.mos)
		if result != tt.expected {
			t.Errorf("GetValuationRating(%v) = %v, expected %v", tt.mos, result, tt.expected)
		}
	}
}

func TestPresentValue(t *testing.T) {
	// $100 in 5 years at 10% discount rate
	pv := PresentValue(100, 10, 5)
	// = 100 / (1.1)^5 = 100 / 1.6105 = 62.09
	expected := 62.09

	if math.Abs(pv-expected) > 0.1 {
		t.Errorf("PV = %v, expected %v", pv, expected)
	}
}

func TestFutureValue(t *testing.T) {
	// $100 today at 10% for 5 years
	fv := FutureValue(100, 10, 5)
	// = 100 * (1.1)^5 = 161.05
	expected := 161.05

	if math.Abs(fv-expected) > 0.1 {
		t.Errorf("FV = %v, expected %v", fv, expected)
	}
}
