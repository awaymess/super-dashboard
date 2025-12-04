package calculations

import "math"

// DCFInputs contains inputs for DCF valuation.
type DCFInputs struct {
	FreeCashFlow       float64 // Current free cash flow
	GrowthRate         float64 // Annual growth rate (percentage)
	TerminalGrowthRate float64 // Terminal growth rate (percentage)
	DiscountRate       float64 // Discount rate / WACC (percentage)
	Years              int     // Projection years
	SharesOutstanding  float64 // Number of shares outstanding
}

// DCFResult contains the result of a DCF valuation.
type DCFResult struct {
	IntrinsicValue           float64
	PerShareValue            float64
	PresentValueOfCashFlows  float64
	TerminalValue            float64
	PresentValueOfTerminal   float64
	ProjectedCashFlows       []CashFlowProjection
}

// CashFlowProjection represents projected cash flow for a year.
type CashFlowProjection struct {
	Year         int
	CashFlow     float64
	PresentValue float64
}

// CalculateDCF performs a Discounted Cash Flow valuation.
func CalculateDCF(inputs DCFInputs) DCFResult {
	projectedCashFlows := make([]CashFlowProjection, 0, inputs.Years)
	totalPVCashFlows := 0.0
	currentCashFlow := inputs.FreeCashFlow

	for year := 1; year <= inputs.Years; year++ {
		currentCashFlow *= (1 + inputs.GrowthRate/100)
		discountFactor := math.Pow(1+inputs.DiscountRate/100, float64(year))
		presentValue := currentCashFlow / discountFactor

		projectedCashFlows = append(projectedCashFlows, CashFlowProjection{
			Year:         year,
			CashFlow:     currentCashFlow,
			PresentValue: presentValue,
		})

		totalPVCashFlows += presentValue
	}

	// Terminal value using Gordon Growth Model
	finalCashFlow := currentCashFlow * (1 + inputs.TerminalGrowthRate/100)
	terminalValue := finalCashFlow / ((inputs.DiscountRate - inputs.TerminalGrowthRate) / 100)
	pvTerminalValue := terminalValue / math.Pow(1+inputs.DiscountRate/100, float64(inputs.Years))

	intrinsicValue := totalPVCashFlows + pvTerminalValue
	perShareValue := 0.0
	if inputs.SharesOutstanding > 0 {
		perShareValue = intrinsicValue / inputs.SharesOutstanding
	}

	return DCFResult{
		IntrinsicValue:          intrinsicValue,
		PerShareValue:           perShareValue,
		PresentValueOfCashFlows: totalPVCashFlows,
		TerminalValue:           terminalValue,
		PresentValueOfTerminal:  pvTerminalValue,
		ProjectedCashFlows:      projectedCashFlows,
	}
}

// CalculateWACC calculates Weighted Average Cost of Capital.
// All inputs as percentages (0-100)
func CalculateWACC(equityWeight, debtWeight, costOfEquity, costOfDebt, taxRate float64) float64 {
	return (equityWeight/100)*(costOfEquity/100) +
		(debtWeight/100)*(costOfDebt/100)*(1-taxRate/100)
}

// CalculateCostOfEquity calculates cost of equity using CAPM.
// Returns cost of equity as percentage
func CalculateCostOfEquity(riskFreeRate, beta, marketReturn float64) float64 {
	return riskFreeRate + beta*(marketReturn-riskFreeRate)
}

// EstimateGrowthRate estimates sustainable growth rate.
// ROE and retention ratio as percentages
func EstimateGrowthRate(roe, retentionRatio float64) float64 {
	return (roe / 100) * (retentionRatio / 100) * 100
}

// ReverseDCF calculates the implied growth rate from current price.
func ReverseDCF(currentPrice, sharesOutstanding, freeCashFlow, discountRate, terminalGrowthRate float64, years int) float64 {
	targetValue := currentPrice * sharesOutstanding

	lowGrowth := 0.0
	highGrowth := 50.0
	midGrowth := 25.0
	tolerance := 0.01
	maxIterations := 100

	for i := 0; i < maxIterations; i++ {
		result := CalculateDCF(DCFInputs{
			FreeCashFlow:       freeCashFlow,
			GrowthRate:         midGrowth,
			TerminalGrowthRate: terminalGrowthRate,
			DiscountRate:       discountRate,
			Years:              years,
			SharesOutstanding:  sharesOutstanding,
		})

		diff := result.IntrinsicValue - targetValue

		if math.Abs(diff) < targetValue*tolerance {
			return midGrowth
		}

		if diff > 0 {
			highGrowth = midGrowth
		} else {
			lowGrowth = midGrowth
		}

		midGrowth = (lowGrowth + highGrowth) / 2
	}

	return midGrowth
}

// CalculateMarginOfSafety calculates the margin of safety for a stock.
func CalculateMarginOfSafety(intrinsicValue, currentPrice float64) float64 {
	if intrinsicValue <= 0 {
		return 0
	}
	return ((intrinsicValue - currentPrice) / intrinsicValue) * 100
}

// ValuationRating represents stock valuation status.
type ValuationRating string

const (
	Undervalued ValuationRating = "undervalued"
	FairValue   ValuationRating = "fair_value"
	Overvalued  ValuationRating = "overvalued"
)

// GetValuationRating determines if a stock is undervalued based on margin of safety.
func GetValuationRating(marginOfSafety float64) ValuationRating {
	switch {
	case marginOfSafety >= 30:
		return Undervalued
	case marginOfSafety >= -10:
		return FairValue
	default:
		return Overvalued
	}
}

// PresentValue calculates the present value of a future cash flow.
func PresentValue(futureValue, discountRate float64, periods int) float64 {
	return futureValue / math.Pow(1+discountRate/100, float64(periods))
}

// FutureValue calculates the future value of a present amount.
func FutureValue(presentValue, growthRate float64, periods int) float64 {
	return presentValue * math.Pow(1+growthRate/100, float64(periods))
}
