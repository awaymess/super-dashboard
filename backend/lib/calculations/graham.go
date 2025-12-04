package calculations

import "math"

// GrahamInputs contains inputs for Graham valuation.
type GrahamInputs struct {
	EPS               float64 // Earnings per share
	BookValuePerShare float64 // Book value per share
	CurrentPrice      float64 // Current stock price
	GrowthRate        float64 // Expected growth rate (optional)
	AAAYield          float64 // AAA corporate bond yield (optional, default 4.4%)
}

// GrahamResult contains the result of Graham valuation.
type GrahamResult struct {
	GrahamNumber        float64
	ModifiedGrahamValue float64
	PELimit             float64
	MarginOfSafety      float64
	Rating              string // "undervalued", "fair_value", "overvalued"
	Analysis            GrahamAnalysis
}

// GrahamAnalysis contains detailed analysis metrics.
type GrahamAnalysis struct {
	PERatio       float64
	PBRatio       float64
	PEG           *float64 // nil if growth rate is 0
	IsDefensive   bool
	IsEnterprising bool
}

// CalculateGrahamNumber computes the classic Graham Number.
// Graham Number = √(22.5 × EPS × Book Value)
func CalculateGrahamNumber(eps, bookValue float64) float64 {
	if eps <= 0 || bookValue <= 0 {
		return 0
	}
	return math.Sqrt(22.5 * eps * bookValue)
}

// CalculateModifiedGrahamValue computes Graham's intrinsic value formula.
// V = EPS × (8.5 + 2g) × 4.4 / Y
// where g is expected growth rate and Y is AAA bond yield
func CalculateModifiedGrahamValue(eps, growthRate, aaaYield float64) float64 {
	if eps <= 0 {
		return 0
	}
	if aaaYield == 0 {
		aaaYield = 4.4 // Default AAA yield
	}

	baseMultiplier := 8.5
	growthMultiplier := 2.0
	expectedReturn := 4.4

	value := (eps * (baseMultiplier + growthMultiplier*growthRate) * expectedReturn) / aaaYield
	return math.Max(0, value)
}

// CalculateGrahamAnalysis performs comprehensive Graham analysis.
func CalculateGrahamAnalysis(inputs GrahamInputs) GrahamResult {
	if inputs.AAAYield == 0 {
		inputs.AAAYield = 4.4
	}

	grahamNumber := CalculateGrahamNumber(inputs.EPS, inputs.BookValuePerShare)
	modifiedGrahamValue := CalculateModifiedGrahamValue(inputs.EPS, inputs.GrowthRate, inputs.AAAYield)

	var peRatio, pbRatio float64
	if inputs.EPS > 0 {
		peRatio = inputs.CurrentPrice / inputs.EPS
	}
	if inputs.BookValuePerShare > 0 {
		pbRatio = inputs.CurrentPrice / inputs.BookValuePerShare
	}

	var peg *float64
	if inputs.GrowthRate > 0 && inputs.EPS > 0 {
		pegValue := peRatio / inputs.GrowthRate
		peg = &pegValue
	}

	// Graham's defensive investor criteria
	peLimit := 15.0
	pbLimit := 1.5
	combinedLimit := 22.5

	isDefensive := peRatio > 0 && peRatio <= peLimit &&
		pbRatio > 0 && pbRatio <= pbLimit &&
		(peRatio*pbRatio) <= combinedLimit

	// Graham's enterprising investor criteria (slightly relaxed)
	isEnterprising := peRatio > 0 && peRatio <= 20 &&
		pbRatio > 0 && pbRatio <= 2

	intrinsicValue := math.Max(grahamNumber, modifiedGrahamValue)
	marginOfSafety := 0.0
	if intrinsicValue > 0 {
		marginOfSafety = ((intrinsicValue - inputs.CurrentPrice) / intrinsicValue) * 100
	}

	var rating string
	switch {
	case marginOfSafety >= 30:
		rating = "undervalued"
	case marginOfSafety >= -10:
		rating = "fair_value"
	default:
		rating = "overvalued"
	}

	return GrahamResult{
		GrahamNumber:        grahamNumber,
		ModifiedGrahamValue: modifiedGrahamValue,
		PELimit:             combinedLimit,
		MarginOfSafety:      marginOfSafety,
		Rating:              rating,
		Analysis: GrahamAnalysis{
			PERatio:       peRatio,
			PBRatio:       pbRatio,
			PEG:           peg,
			IsDefensive:   isDefensive,
			IsEnterprising: isEnterprising,
		},
	}
}

// ScreenDefensiveStocks filters stocks that meet Graham's defensive criteria.
type StockForScreening struct {
	Symbol    string
	EPS       float64
	BookValue float64
	Price     float64
}

// ScreenDefensiveStocks returns symbols of stocks meeting defensive criteria.
func ScreenDefensiveStocks(stocks []StockForScreening) []string {
	result := make([]string, 0)
	for _, stock := range stocks {
		if stock.EPS <= 0 || stock.BookValue <= 0 {
			continue
		}
		pe := stock.Price / stock.EPS
		pb := stock.Price / stock.BookValue
		if pe <= 15 && pb <= 1.5 && (pe*pb) <= 22.5 {
			result = append(result, stock.Symbol)
		}
	}
	return result
}

// NCAVResult contains Net Current Asset Value calculation.
type NCAVResult struct {
	NCAV         float64
	NCAVPerShare float64
}

// CalculateNCAV computes Net Current Asset Value.
// NCAV = Current Assets - Total Liabilities
func CalculateNCAV(currentAssets, totalLiabilities, sharesOutstanding float64) NCAVResult {
	ncav := currentAssets - totalLiabilities
	ncavPerShare := 0.0
	if sharesOutstanding > 0 {
		ncavPerShare = ncav / sharesOutstanding
	}
	return NCAVResult{
		NCAV:         ncav,
		NCAVPerShare: ncavPerShare,
	}
}

// IsNetNet checks if a stock is a Graham "net-net" bargain.
// Net-net: Market cap < 2/3 of NCAV
func IsNetNet(currentAssets, totalLiabilities, marketCap float64) bool {
	ncav := currentAssets - totalLiabilities
	return marketCap < ncav*0.67
}

// BuffettIntrinsicValue calculates intrinsic value using Buffett's owner earnings method.
// Owner Earnings = Net Income + Depreciation - CapEx - Working Capital Changes
func BuffettIntrinsicValue(ownerEarnings, growthRate, discountRate float64, years int) float64 {
	if discountRate <= growthRate {
		return 0 // Invalid inputs
	}

	totalPV := 0.0
	currentEarnings := ownerEarnings

	for year := 1; year <= years; year++ {
		currentEarnings *= (1 + growthRate/100)
		pv := currentEarnings / math.Pow(1+discountRate/100, float64(year))
		totalPV += pv
	}

	// Terminal value
	terminalEarnings := currentEarnings * (1 + 3.0/100) // 3% perpetual growth
	terminalValue := terminalEarnings / ((discountRate - 3) / 100)
	pvTerminal := terminalValue / math.Pow(1+discountRate/100, float64(years))

	return totalPV + pvTerminal
}

// PBVValuation calculates fair price using P/BV method.
// Commonly used for banks and real estate companies.
func PBVValuation(bookValue, fairPBRatio float64) float64 {
	return bookValue * fairPBRatio
}

// PEValuation calculates fair price using P/E method.
func PEValuation(eps, fairPERatio float64) float64 {
	return eps * fairPERatio
}
