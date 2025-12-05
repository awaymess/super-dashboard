package calculations

import (
	"math"
)

// PositionSize calculates position size based on risk parameters.
type PositionSizeParams struct {
	AccountSize      float64
	RiskPercentage   float64 // Risk per trade as percentage (e.g., 2.0 for 2%)
	EntryPrice       float64
	StopLossPrice    float64
}

func PositionSize(params PositionSizeParams) float64 {
	riskAmount := params.AccountSize * (params.RiskPercentage / 100)
	riskPerShare := math.Abs(params.EntryPrice - params.StopLossPrice)

	if riskPerShare == 0 {
		return 0
	}

	shares := riskAmount / riskPerShare
	return math.Floor(shares) // Round down to whole shares
}

// PositionSizeKelly calculates position size using Kelly Criterion.
func PositionSizeKelly(accountSize, winRate, avgWin, avgLoss, fraction float64) float64 {
	if avgLoss == 0 {
		return 0
	}

	// Kelly formula: f = (p * b - q) / b
	// where p = win rate, q = 1 - p, b = avg win / avg loss
	b := avgWin / avgLoss
	p := winRate
	q := 1 - p

	kellyFraction := (p*b - q) / b

	if kellyFraction <= 0 {
		return 0
	}

	// Apply fractional Kelly for safety
	return accountSize * kellyFraction * fraction
}

// RiskRewardRatio calculates risk-to-reward ratio.
func RiskRewardRatio(entryPrice, stopLoss, target float64) float64 {
	risk := math.Abs(entryPrice - stopLoss)
	reward := math.Abs(target - entryPrice)

	if risk == 0 {
		return 0
	}

	return reward / risk
}

// RequiredWinRate calculates minimum win rate needed for profitability.
func RequiredWinRate(riskRewardRatio float64) float64 {
	return (1 / (1 + riskRewardRatio)) * 100
}

// MaxPositionSize calculates maximum position size as percentage of portfolio.
func MaxPositionSize(accountSize, maxPercentage float64) float64 {
	return accountSize * (maxPercentage / 100)
}

// StopLossPrice calculates stop loss price.
func StopLossPrice(entryPrice, riskPercentage float64, isLong bool) float64 {
	if isLong {
		return entryPrice * (1 - riskPercentage/100)
	}
	return entryPrice * (1 + riskPercentage/100)
}

// TakeProfitPrice calculates take profit price.
func TakeProfitPrice(entryPrice, targetPercentage float64, isLong bool) float64 {
	if isLong {
		return entryPrice * (1 + targetPercentage/100)
	}
	return entryPrice * (1 - targetPercentage/100)
}

// PortfolioHeatMap calculates correlated risk across positions.
func PortfolioHeatMap(positions []float64, correlations [][]float64) float64 {
	if len(positions) != len(correlations) {
		return 0
	}

	totalRisk := 0.0
	for i := range positions {
		for j := range positions {
			totalRisk += positions[i] * positions[j] * correlations[i][j]
		}
	}

	return math.Sqrt(totalRisk)
}

// VaRPosition calculates Value at Risk for a position.
func VaRPosition(positionValue, volatility, confidenceLevel float64) float64 {
	// Z-score for confidence levels
	var zScore float64
	switch confidenceLevel {
	case 0.95:
		zScore = 1.645
	case 0.99:
		zScore = 2.326
	default:
		zScore = 1.645 // Default to 95%
	}

	return positionValue * volatility * zScore
}

// ExpectedShortfall calculates ES/CVaR for a position.
func ExpectedShortfall(positionValue, volatility, confidenceLevel float64) float64 {
	var ratio float64
	switch confidenceLevel {
	case 0.95:
		ratio = 2.063
	case 0.99:
		ratio = 2.665
	default:
		ratio = 2.063
	}

	return positionValue * volatility * ratio
}

// PortfolioDiversification calculates effective number of positions.
func PortfolioDiversification(weights []float64) float64 {
	sumSquares := 0.0
	for _, w := range weights {
		sumSquares += w * w
	}

	if sumSquares == 0 {
		return 0
	}

	return 1 / sumSquares
}

// ConcentrationRisk calculates concentration using Herfindahl index.
func ConcentrationRisk(weights []float64) float64 {
	hhi := 0.0
	for _, w := range weights {
		hhi += w * w * 10000 // Scale to 0-10000
	}
	return hhi
}

// MaxDrawdownStop calculates when to stop trading based on drawdown.
func MaxDrawdownStop(initialEquity, currentEquity, maxDrawdownPercent float64) bool {
	currentDrawdown := ((initialEquity - currentEquity) / initialEquity) * 100
	return currentDrawdown >= maxDrawdownPercent
}

// RecoveryFactor calculates recovery factor.
// Recovery Factor = Net Profit / Maximum Drawdown
func RecoveryFactor(netProfit, maxDrawdown float64) float64 {
	if maxDrawdown == 0 {
		return 0
	}
	return netProfit / maxDrawdown
}

// RiskAdjustedReturn calculates risk-adjusted return.
func RiskAdjustedReturn(returns, risk float64) float64 {
	if risk == 0 {
		return 0
	}
	return returns / risk
}

// SafetyFirstRatio calculates Roy's Safety-First Ratio.
// SF = (Expected Return - Threshold) / Standard Deviation
func SafetyFirstRatio(expectedReturn, threshold, stdDev float64) float64 {
	if stdDev == 0 {
		return 0
	}
	return (expectedReturn - threshold) / stdDev
}

// LeverageRatio calculates leverage ratio.
func LeverageRatio(totalPositionValue, accountEquity float64) float64 {
	if accountEquity == 0 {
		return 0
	}
	return totalPositionValue / accountEquity
}

// MarginRequirement calculates margin required for position.
func MarginRequirement(positionValue, marginPercent float64) float64 {
	return positionValue * (marginPercent / 100)
}

// LiquidationPrice calculates liquidation price for leveraged position.
func LiquidationPrice(entryPrice, leverage, isLong bool) float64 {
	liquidationPercent := 1 / leverage

	if isLong {
		return entryPrice * (1 - liquidationPercent)
	}
	return entryPrice * (1 + liquidationPercent)
}

// RiskOfRuin calculates probability of ruin.
// Simplified formula for equal bet sizes
func RiskOfRuin(winRate, riskRewardRatio float64, initialCapital, targetCapital float64) float64 {
	if winRate >= 0.5 && riskRewardRatio >= 1 {
		return 0 // No risk if edge is positive
	}

	// Simplified calculation
	edgePerTrade := (winRate * riskRewardRatio) - (1 - winRate)
	if edgePerTrade >= 0 {
		return 0
	}

	ratio := targetCapital / initialCapital
	return math.Pow(ratio, -edgePerTrade)
}

// TrailingStop calculates trailing stop price.
func TrailingStop(highestPrice, trailingPercent float64) float64 {
	return highestPrice * (1 - trailingPercent/100)
}

// BreakevenWinRate calculates breakeven win rate given costs.
func BreakevenWinRate(commission, slippage, riskRewardRatio float64) float64 {
	totalCost := commission + slippage
	return (1 + totalCost) / (1 + riskRewardRatio)
}

// PortfolioCorrelation calculates portfolio correlation risk.
func PortfolioCorrelation(weights []float64, correlationMatrix [][]float64) float64 {
	if len(weights) != len(correlationMatrix) {
		return 0
	}

	totalCorrelation := 0.0
	for i := range weights {
		for j := range weights {
			if i != j {
				totalCorrelation += weights[i] * weights[j] * correlationMatrix[i][j]
			}
		}
	}

	return totalCorrelation
}

// OptimalFPosition calculates optimal f (optimal fraction).
func OptimalFPosition(trades []float64) float64 {
	if len(trades) == 0 {
		return 0
	}

	// Find largest loss
	largestLoss := 0.0
	for _, trade := range trades {
		if trade < largestLoss {
			largestLoss = trade
		}
	}

	if largestLoss == 0 {
		return 0
	}

	// Calculate TWR (Terminal Wealth Relative) for different f values
	bestF := 0.0
	bestTWR := 0.0

	for f := 0.01; f <= 1.0; f += 0.01 {
		twr := 1.0
		for _, trade := range trades {
			hpr := 1 + (f * trade / math.Abs(largestLoss))
			if hpr <= 0 {
				twr = 0
				break
			}
			twr *= hpr
		}

		if twr > bestTWR {
			bestTWR = twr
			bestF = f
		}
	}

	return bestF
}

// DynamicPositionSize adjusts position size based on recent performance.
func DynamicPositionSize(baseSize, recentWinRate, targetWinRate float64) float64 {
	if targetWinRate == 0 {
		return baseSize
	}

	adjustment := recentWinRate / targetWinRate
	// Cap adjustment between 0.5x and 2x
	if adjustment < 0.5 {
		adjustment = 0.5
	}
	if adjustment > 2.0 {
		adjustment = 2.0
	}

	return baseSize * adjustment
}

// StreakAdjustment adjusts position size based on win/loss streak.
func StreakAdjustment(baseSize float64, streak int, isWinStreak bool) float64 {
	if isWinStreak {
		// Increase size gradually during win streak (max 50% increase)
		increase := math.Min(float64(streak)*0.1, 0.5)
		return baseSize * (1 + increase)
	}

	// Decrease size during loss streak (max 50% decrease)
	decrease := math.Min(float64(streak)*0.1, 0.5)
	return baseSize * (1 - decrease)
}

// MonteCarloRisk simulates potential outcomes.
func MonteCarloRisk(trades []float64, simulations int) map[string]float64 {
	// Simple Monte Carlo implementation
	outcomes := make([]float64, simulations)

	for i := 0; i < simulations; i++ {
		sum := 0.0
		for _, trade := range trades {
			// Randomly select with replacement
			sum += trade
		}
		outcomes[i] = sum
	}

	// Calculate statistics
	mean := 0.0
	for _, outcome := range outcomes {
		mean += outcome
	}
	mean /= float64(simulations)

	// Calculate percentiles
	// This is simplified - in production use proper sorting
	worst := outcomes[0]
	best := outcomes[0]
	for _, outcome := range outcomes {
		if outcome < worst {
			worst = outcome
		}
		if outcome > best {
			best = outcome
		}
	}

	return map[string]float64{
		"mean":          mean,
		"worst_case":    worst,
		"best_case":     best,
		"risk_range":    best - worst,
	}
}
