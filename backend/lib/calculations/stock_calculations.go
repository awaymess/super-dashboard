package calculations

import (
	"math"
)

// DCFValuation calculates stock fair value using Discounted Cash Flow.
type DCFParams struct {
	FreeCashFlow       float64
	GrowthRate         float64 // Expected growth rate (e.g., 0.10 for 10%)
	TerminalGrowthRate float64 // Terminal growth rate (typically 0.02-0.03)
	DiscountRate       float64 // WACC (Weighted Average Cost of Capital)
	Years              int     // Projection period (typically 5-10 years)
	SharesOutstanding  float64
}

func DCFValuation(params DCFParams) float64 {
	presentValue := 0.0
	currentCF := params.FreeCashFlow

	// Calculate present value of projected cash flows
	for i := 1; i <= params.Years; i++ {
		currentCF *= (1 + params.GrowthRate)
		discountFactor := math.Pow(1+params.DiscountRate, float64(i))
		presentValue += currentCF / discountFactor
	}

	// Calculate terminal value using perpetuity growth model
	terminalCF := currentCF * (1 + params.TerminalGrowthRate)
	terminalValue := terminalCF / (params.DiscountRate - params.TerminalGrowthRate)
	discountedTerminalValue := terminalValue / math.Pow(1+params.DiscountRate, float64(params.Years))

	// Total enterprise value
	enterpriseValue := presentValue + discountedTerminalValue

	// Fair value per share
	fairValuePerShare := enterpriseValue / params.SharesOutstanding

	return fairValuePerShare
}

// GrahamNumber calculates Benjamin Graham's intrinsic value.
// Graham Number = sqrt(22.5 × EPS × Book Value Per Share)
func GrahamNumber(eps, bookValuePerShare float64) float64 {
	if eps <= 0 || bookValuePerShare <= 0 {
		return 0
	}
	return math.Sqrt(22.5 * eps * bookValuePerShare)
}

// GrahamIntrinsicValue calculates intrinsic value using Graham's formula.
// IV = EPS × (8.5 + 2g) × 4.4 / Y
// where g = expected growth rate, Y = current yield on AAA bonds
func GrahamIntrinsicValue(eps, growthRate, bondYield float64) float64 {
	if bondYield == 0 {
		bondYield = 0.045 // Default 4.5%
	}
	return eps * (8.5 + 2*growthRate*100) * 4.4 / (bondYield * 100)
}

// PEValuation calculates fair value using P/E ratio.
func PEValuation(eps, targetPE float64) float64 {
	return eps * targetPE
}

// PEGRatio calculates Price/Earnings to Growth ratio.
// PEG = (Price / EPS) / Growth Rate
func PEGRatio(price, eps, growthRate float64) float64 {
	if growthRate == 0 || eps == 0 {
		return 0
	}
	pe := price / eps
	return pe / (growthRate * 100)
}

// DividendDiscountModel calculates stock value using DDM.
// Value = D1 / (r - g)
// where D1 = next year dividend, r = required return, g = growth rate
func DividendDiscountModel(currentDividend, growthRate, requiredReturn float64) float64 {
	if requiredReturn <= growthRate {
		return 0 // Invalid: growth rate must be less than required return
	}
	nextDividend := currentDividend * (1 + growthRate)
	return nextDividend / (requiredReturn - growthRate)
}

// PriceToBook calculates P/B ratio.
func PriceToBook(price, bookValuePerShare float64) float64 {
	if bookValuePerShare == 0 {
		return 0
	}
	return price / bookValuePerShare
}

// PriceToSales calculates P/S ratio.
func PriceToSales(marketCap, revenue float64) float64 {
	if revenue == 0 {
		return 0
	}
	return marketCap / revenue
}

// EnterpriseValue calculates enterprise value.
// EV = Market Cap + Debt - Cash
func EnterpriseValue(marketCap, totalDebt, cash float64) float64 {
	return marketCap + totalDebt - cash
}

// EVToEBITDA calculates EV/EBITDA multiple.
func EVToEBITDA(enterpriseValue, ebitda float64) float64 {
	if ebitda == 0 {
		return 0
	}
	return enterpriseValue / ebitda
}

// EVToSales calculates EV/Sales multiple.
func EVToSales(enterpriseValue, revenue float64) float64 {
	if revenue == 0 {
		return 0
	}
	return enterpriseValue / revenue
}

// DebtToEquity calculates debt-to-equity ratio.
func DebtToEquity(totalDebt, totalEquity float64) float64 {
	if totalEquity == 0 {
		return 0
	}
	return totalDebt / totalEquity
}

// CurrentRatio calculates current ratio (liquidity).
func CurrentRatio(currentAssets, currentLiabilities float64) float64 {
	if currentLiabilities == 0 {
		return 0
	}
	return currentAssets / currentLiabilities
}

// QuickRatio calculates quick ratio (acid-test ratio).
// Quick Ratio = (Current Assets - Inventory) / Current Liabilities
func QuickRatio(currentAssets, inventory, currentLiabilities float64) float64 {
	if currentLiabilities == 0 {
		return 0
	}
	return (currentAssets - inventory) / currentLiabilities
}

// ROE calculates Return on Equity.
func ROE(netIncome, shareholderEquity float64) float64 {
	if shareholderEquity == 0 {
		return 0
	}
	return (netIncome / shareholderEquity) * 100
}

// ROA calculates Return on Assets.
func ROA(netIncome, totalAssets float64) float64 {
	if totalAssets == 0 {
		return 0
	}
	return (netIncome / totalAssets) * 100
}

// ROIC calculates Return on Invested Capital.
// ROIC = NOPAT / Invested Capital
func ROIC(nopat, investedCapital float64) float64 {
	if investedCapital == 0 {
		return 0
	}
	return (nopat / investedCapital) * 100
}

// EarningsYield calculates earnings yield (inverse of P/E).
func EarningsYield(eps, price float64) float64 {
	if price == 0 {
		return 0
	}
	return (eps / price) * 100
}

// DividendYield calculates dividend yield.
func DividendYield(annualDividend, price float64) float64 {
	if price == 0 {
		return 0
	}
	return (annualDividend / price) * 100
}

// PayoutRatio calculates dividend payout ratio.
func PayoutRatio(dividendPerShare, earningsPerShare float64) float64 {
	if earningsPerShare == 0 {
		return 0
	}
	return (dividendPerShare / earningsPerShare) * 100
}

// RetentionRatio calculates earnings retention ratio.
func RetentionRatio(dividendPerShare, earningsPerShare float64) float64 {
	return 100 - PayoutRatio(dividendPerShare, earningsPerShare)
}

// OperatingMargin calculates operating profit margin.
func OperatingMargin(operatingIncome, revenue float64) float64 {
	if revenue == 0 {
		return 0
	}
	return (operatingIncome / revenue) * 100
}

// NetProfitMargin calculates net profit margin.
func NetProfitMargin(netIncome, revenue float64) float64 {
	if revenue == 0 {
		return 0
	}
	return (netIncome / revenue) * 100
}

// GrossMargin calculates gross profit margin.
func GrossMargin(revenue, cogs float64) float64 {
	if revenue == 0 {
		return 0
	}
	return ((revenue - cogs) / revenue) * 100
}

// AssetTurnover calculates asset turnover ratio.
func AssetTurnover(revenue, totalAssets float64) float64 {
	if totalAssets == 0 {
		return 0
	}
	return revenue / totalAssets
}

// InventoryTurnover calculates inventory turnover ratio.
func InventoryTurnover(cogs, averageInventory float64) float64 {
	if averageInventory == 0 {
		return 0
	}
	return cogs / averageInventory
}

// ReceivablesTurnover calculates receivables turnover ratio.
func ReceivablesTurnover(revenue, averageReceivables float64) float64 {
	if averageReceivables == 0 {
		return 0
	}
	return revenue / averageReceivables
}

// DaysInventoryOutstanding calculates DIO.
func DaysInventoryOutstanding(cogs, averageInventory float64) float64 {
	turnover := InventoryTurnover(cogs, averageInventory)
	if turnover == 0 {
		return 0
	}
	return 365 / turnover
}

// DaysSalesOutstanding calculates DSO.
func DaysSalesOutstanding(revenue, averageReceivables float64) float64 {
	turnover := ReceivablesTurnover(revenue, averageReceivables)
	if turnover == 0 {
		return 0
	}
	return 365 / turnover
}

// CashConversionCycle calculates CCC.
// CCC = DIO + DSO - DPO
func CashConversionCycle(dio, dso, dpo float64) float64 {
	return dio + dso - dpo
}

// WACC calculates Weighted Average Cost of Capital.
type WACCParams struct {
	MarketValueEquity float64
	MarketValueDebt   float64
	CostOfEquity      float64
	CostOfDebt        float64
	TaxRate           float64
}

func WACC(params WACCParams) float64 {
	totalValue := params.MarketValueEquity + params.MarketValueDebt

	if totalValue == 0 {
		return 0
	}

	equityWeight := params.MarketValueEquity / totalValue
	debtWeight := params.MarketValueDebt / totalValue

	wacc := (equityWeight * params.CostOfEquity) +
		(debtWeight * params.CostOfDebt * (1 - params.TaxRate))

	return wacc * 100
}

// CAPM calculates Cost of Equity using Capital Asset Pricing Model.
// Cost of Equity = Risk-Free Rate + Beta × (Market Return - Risk-Free Rate)
func CAPM(riskFreeRate, beta, marketReturn float64) float64 {
	return riskFreeRate + beta*(marketReturn-riskFreeRate)
}

// AltmanZScore calculates Altman Z-Score for bankruptcy prediction.
func AltmanZScore(workingCapital, retainedEarnings, ebit, marketValueEquity, totalAssets, totalLiabilities, sales float64) float64 {
	x1 := workingCapital / totalAssets
	x2 := retainedEarnings / totalAssets
	x3 := ebit / totalAssets
	x4 := marketValueEquity / totalLiabilities
	x5 := sales / totalAssets

	z := 1.2*x1 + 1.4*x2 + 3.3*x3 + 0.6*x4 + 1.0*x5

	return z
}

// PiotroskiFScore calculates Piotroski F-Score (fundamental strength).
func PiotroskiFScore(params map[string]float64) int {
	score := 0

	// Profitability (4 points)
	if params["netIncome"] > 0 {
		score++
	}
	if params["roa"] > 0 {
		score++
	}
	if params["operatingCashFlow"] > 0 {
		score++
	}
	if params["operatingCashFlow"] > params["netIncome"] {
		score++
	}

	// Leverage, Liquidity, Source of Funds (3 points)
	if params["currentDebtToEquity"] < params["previousDebtToEquity"] {
		score++
	}
	if params["currentRatio"] > params["previousCurrentRatio"] {
		score++
	}
	if params["newSharesIssued"] == 0 {
		score++
	}

	// Operating Efficiency (2 points)
	if params["currentGrossMargin"] > params["previousGrossMargin"] {
		score++
	}
	if params["currentAssetTurnover"] > params["previousAssetTurnover"] {
		score++
	}

	return score
}

// IntrinsicValueMargin calculates margin of safety.
func IntrinsicValueMargin(intrinsicValue, currentPrice float64) float64 {
	if intrinsicValue == 0 {
		return 0
	}
	return ((intrinsicValue - currentPrice) / intrinsicValue) * 100
}

// TargetPrice calculates target price based on target P/E and expected EPS.
func TargetPrice(expectedEPS, targetPE float64) float64 {
	return expectedEPS * targetPE
}

// UpsidePotential calculates upside potential percentage.
func UpsidePotential(targetPrice, currentPrice float64) float64 {
	if currentPrice == 0 {
		return 0
	}
	return ((targetPrice - currentPrice) / currentPrice) * 100
}
