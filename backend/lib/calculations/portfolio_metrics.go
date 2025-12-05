package calculations

import (
	"math"
	"sort"
)

// PortfolioReturn calculates total return of a portfolio.
func PortfolioReturn(initialValue, currentValue float64) float64 {
	if initialValue == 0 {
		return 0
	}
	return ((currentValue - initialValue) / initialValue) * 100
}

// SharpeRatio calculates Sharpe Ratio.
// Sharpe = (Portfolio Return - Risk-Free Rate) / Standard Deviation
func SharpeRatio(returns []float64, riskFreeRate float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	avgReturn := 0.0
	for _, r := range returns {
		avgReturn += r
	}
	avgReturn /= float64(len(returns))

	// Calculate standard deviation
	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-avgReturn, 2)
	}
	variance /= float64(len(returns))
	stdDev := math.Sqrt(variance)

	if stdDev == 0 {
		return 0
	}

	return (avgReturn - riskFreeRate) / stdDev
}

// SortinoRatio calculates Sortino Ratio (only considers downside deviation).
func SortinoRatio(returns []float64, targetReturn float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	avgReturn := 0.0
	for _, r := range returns {
		avgReturn += r
	}
	avgReturn /= float64(len(returns))

	// Calculate downside deviation
	downsideVariance := 0.0
	count := 0
	for _, r := range returns {
		if r < targetReturn {
			downsideVariance += math.Pow(r-targetReturn, 2)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	downsideVariance /= float64(count)
	downsideDeviation := math.Sqrt(downsideVariance)

	if downsideDeviation == 0 {
		return 0
	}

	return (avgReturn - targetReturn) / downsideDeviation
}

// MaxDrawdown calculates maximum drawdown percentage.
func MaxDrawdown(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	maxDrawdown := 0.0
	peak := values[0]

	for _, value := range values {
		if value > peak {
			peak = value
		}

		drawdown := ((peak - value) / peak) * 100
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}

// CalmarRatio calculates Calmar Ratio.
// Calmar = Annual Return / Maximum Drawdown
func CalmarRatio(annualReturn, maxDrawdown float64) float64 {
	if maxDrawdown == 0 {
		return 0
	}
	return annualReturn / maxDrawdown
}

// ValueAtRisk calculates Value at Risk at given confidence level.
// VaR = percentile of returns distribution
func ValueAtRisk(returns []float64, confidenceLevel float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	sorted := make([]float64, len(returns))
	copy(sorted, returns)
	sort.Float64s(sorted)

	index := int(float64(len(sorted)) * (1 - confidenceLevel))
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// ConditionalValueAtRisk calculates CVaR (Expected Shortfall).
// CVaR = average of returns below VaR threshold
func ConditionalValueAtRisk(returns []float64, confidenceLevel float64) float64 {
	var := ValueAtRisk(returns, confidenceLevel)

	sum := 0.0
	count := 0
	for _, r := range returns {
		if r <= var {
			sum += r
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

// Beta calculates portfolio beta (systematic risk).
func Beta(portfolioReturns, marketReturns []float64) float64 {
	if len(portfolioReturns) != len(marketReturns) || len(portfolioReturns) == 0 {
		return 1.0 // Default beta
	}

	// Calculate means
	portfolioMean := 0.0
	marketMean := 0.0
	for i := range portfolioReturns {
		portfolioMean += portfolioReturns[i]
		marketMean += marketReturns[i]
	}
	portfolioMean /= float64(len(portfolioReturns))
	marketMean /= float64(len(marketReturns))

	// Calculate covariance and variance
	covariance := 0.0
	marketVariance := 0.0
	for i := range portfolioReturns {
		covariance += (portfolioReturns[i] - portfolioMean) * (marketReturns[i] - marketMean)
		marketVariance += math.Pow(marketReturns[i]-marketMean, 2)
	}

	if marketVariance == 0 {
		return 1.0
	}

	return covariance / marketVariance
}

// Alpha calculates Jensen's Alpha.
// Alpha = Portfolio Return - (Risk-Free Rate + Beta × (Market Return - Risk-Free Rate))
func Alpha(portfolioReturn, marketReturn, riskFreeRate, beta float64) float64 {
	expectedReturn := riskFreeRate + beta*(marketReturn-riskFreeRate)
	return portfolioReturn - expectedReturn
}

// TreynorRatio calculates Treynor Ratio.
// Treynor = (Portfolio Return - Risk-Free Rate) / Beta
func TreynorRatio(portfolioReturn, riskFreeRate, beta float64) float64 {
	if beta == 0 {
		return 0
	}
	return (portfolioReturn - riskFreeRate) / beta
}

// InformationRatio calculates Information Ratio.
// IR = (Portfolio Return - Benchmark Return) / Tracking Error
func InformationRatio(portfolioReturns, benchmarkReturns []float64) float64 {
	if len(portfolioReturns) != len(benchmarkReturns) || len(portfolioReturns) == 0 {
		return 0
	}

	// Calculate excess returns
	excessReturns := make([]float64, len(portfolioReturns))
	for i := range portfolioReturns {
		excessReturns[i] = portfolioReturns[i] - benchmarkReturns[i]
	}

	// Calculate average excess return
	avgExcessReturn := 0.0
	for _, er := range excessReturns {
		avgExcessReturn += er
	}
	avgExcessReturn /= float64(len(excessReturns))

	// Calculate tracking error (standard deviation of excess returns)
	trackingErrorVariance := 0.0
	for _, er := range excessReturns {
		trackingErrorVariance += math.Pow(er-avgExcessReturn, 2)
	}
	trackingErrorVariance /= float64(len(excessReturns))
	trackingError := math.Sqrt(trackingErrorVariance)

	if trackingError == 0 {
		return 0
	}

	return avgExcessReturn / trackingError
}

// Correlation calculates correlation between two return series.
func Correlation(returns1, returns2 []float64) float64 {
	if len(returns1) != len(returns2) || len(returns1) == 0 {
		return 0
	}

	mean1 := 0.0
	mean2 := 0.0
	for i := range returns1 {
		mean1 += returns1[i]
		mean2 += returns2[i]
	}
	mean1 /= float64(len(returns1))
	mean2 /= float64(len(returns2))

	numerator := 0.0
	variance1 := 0.0
	variance2 := 0.0

	for i := range returns1 {
		diff1 := returns1[i] - mean1
		diff2 := returns2[i] - mean2
		numerator += diff1 * diff2
		variance1 += diff1 * diff1
		variance2 += diff2 * diff2
	}

	denominator := math.Sqrt(variance1 * variance2)
	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// PortfolioVolatility calculates annualized volatility.
func PortfolioVolatility(returns []float64, periodsPerYear int) float64 {
	if len(returns) == 0 {
		return 0
	}

	mean := 0.0
	for _, r := range returns {
		mean += r
	}
	mean /= float64(len(returns))

	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-mean, 2)
	}
	variance /= float64(len(returns))

	// Annualize
	return math.Sqrt(variance * float64(periodsPerYear))
}

// DownsideDeviation calculates downside deviation (semi-deviation).
func DownsideDeviation(returns []float64, targetReturn float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	downsideVariance := 0.0
	count := 0

	for _, r := range returns {
		if r < targetReturn {
			downsideVariance += math.Pow(r-targetReturn, 2)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return math.Sqrt(downsideVariance / float64(count))
}

// UpsideDeviation calculates upside deviation.
func UpsideDeviation(returns []float64, targetReturn float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	upsideVariance := 0.0
	count := 0

	for _, r := range returns {
		if r > targetReturn {
			upsideVariance += math.Pow(r-targetReturn, 2)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return math.Sqrt(upsideVariance / float64(count))
}

// OmegaRatio calculates Omega Ratio.
// Omega = Probability-weighted gains / Probability-weighted losses
func OmegaRatio(returns []float64, threshold float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	gains := 0.0
	losses := 0.0

	for _, r := range returns {
		if r > threshold {
			gains += r - threshold
		} else {
			losses += threshold - r
		}
	}

	if losses == 0 {
		return 0
	}

	return gains / losses
}

// DrawdownDuration calculates duration of drawdown periods.
type DrawdownPeriod struct {
	StartIndex int
	EndIndex   int
	Duration   int
	Depth      float64
}

func DrawdownDurations(values []float64) []DrawdownPeriod {
	if len(values) == 0 {
		return []DrawdownPeriod{}
	}

	periods := []DrawdownPeriod{}
	peak := values[0]
	peakIndex := 0
	inDrawdown := false
	var currentPeriod DrawdownPeriod

	for i, value := range values {
		if value > peak {
			if inDrawdown {
				currentPeriod.EndIndex = i - 1
				currentPeriod.Duration = currentPeriod.EndIndex - currentPeriod.StartIndex + 1
				periods = append(periods, currentPeriod)
				inDrawdown = false
			}
			peak = value
			peakIndex = i
		} else if value < peak {
			if !inDrawdown {
				inDrawdown = true
				currentPeriod = DrawdownPeriod{
					StartIndex: peakIndex,
				}
			}
			drawdown := ((peak - value) / peak) * 100
			if drawdown > currentPeriod.Depth {
				currentPeriod.Depth = drawdown
			}
		}
	}

	// Close last drawdown if still ongoing
	if inDrawdown {
		currentPeriod.EndIndex = len(values) - 1
		currentPeriod.Duration = currentPeriod.EndIndex - currentPeriod.StartIndex + 1
		periods = append(periods, currentPeriod)
	}

	return periods
}

// RecoveryTime calculates time to recover from drawdown.
func RecoveryTime(values []float64) int {
	drawdowns := DrawdownDurations(values)
	if len(drawdowns) == 0 {
		return 0
	}

	// Return duration of longest drawdown
	maxDuration := 0
	for _, dd := range drawdowns {
		if dd.Duration > maxDuration {
			maxDuration = dd.Duration
		}
	}

	return maxDuration
}

// WinRate calculates percentage of winning periods.
func WinRate(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	wins := 0
	for _, r := range returns {
		if r > 0 {
			wins++
		}
	}

	return (float64(wins) / float64(len(returns))) * 100
}

// ProfitFactor calculates profit factor.
// Profit Factor = Gross Profit / Gross Loss
func ProfitFactor(returns []float64) float64 {
	grossProfit := 0.0
	grossLoss := 0.0

	for _, r := range returns {
		if r > 0 {
			grossProfit += r
		} else {
			grossLoss += math.Abs(r)
		}
	}

	if grossLoss == 0 {
		return 0
	}

	return grossProfit / grossLoss
}

// ExpectancyRatio calculates expectancy (average win × win rate - average loss × loss rate).
func ExpectancyRatio(returns []float64) float64 {
	if len(returns) == 0 {
		return 0
	}

	wins := 0.0
	losses := 0.0
	winCount := 0
	lossCount := 0

	for _, r := range returns {
		if r > 0 {
			wins += r
			winCount++
		} else {
			losses += math.Abs(r)
			lossCount++
		}
	}

	avgWin := 0.0
	avgLoss := 0.0
	winRate := 0.0
	lossRate := 0.0

	if winCount > 0 {
		avgWin = wins / float64(winCount)
		winRate = float64(winCount) / float64(len(returns))
	}

	if lossCount > 0 {
		avgLoss = losses / float64(lossCount)
		lossRate = float64(lossCount) / float64(len(returns))
	}

	return (avgWin * winRate) - (avgLoss * lossRate)
}
