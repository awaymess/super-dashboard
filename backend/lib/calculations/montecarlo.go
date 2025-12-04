package calculations

import (
	"math"
	"math/rand"
	"sort"
)

// MonteCarloResult contains Monte Carlo simulation results.
type MonteCarloResult struct {
	Mean               float64
	Median             float64
	StandardDeviation  float64
	Percentile5        float64
	Percentile25       float64
	Percentile75       float64
	Percentile95       float64
	MinValue           float64
	MaxValue           float64
	ProbabilityOfLoss  float64
	ProbabilityOfGain  float64
	Simulations        int
	Distribution       []float64 // Optional: sample of simulated values
}

// MonteCarloConfig contains configuration for Monte Carlo simulation.
type MonteCarloConfig struct {
	Simulations      int     // Number of simulations to run
	InitialValue     float64 // Starting value (e.g., portfolio value)
	ExpectedReturn   float64 // Expected annual return (percentage)
	Volatility       float64 // Annual volatility / std dev (percentage)
	TimeHorizonYears float64 // Investment horizon in years
	Seed             int64   // Random seed (0 for random)
}

// RunMonteCarloSimulation performs Monte Carlo simulation for investment returns.
func RunMonteCarloSimulation(config MonteCarloConfig) MonteCarloResult {
	if config.Simulations == 0 {
		config.Simulations = 10000
	}
	if config.InitialValue == 0 {
		config.InitialValue = 100000
	}
	if config.TimeHorizonYears == 0 {
		config.TimeHorizonYears = 1
	}

	// Create random source
	var rng *rand.Rand
	if config.Seed != 0 {
		rng = rand.New(rand.NewSource(config.Seed))
	} else {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	results := make([]float64, config.Simulations)
	mu := config.ExpectedReturn / 100    // Convert to decimal
	sigma := config.Volatility / 100     // Convert to decimal
	dt := config.TimeHorizonYears

	lossCount := 0

	for i := 0; i < config.Simulations; i++ {
		// Geometric Brownian Motion
		// S(t) = S(0) * exp((mu - sigma^2/2)*t + sigma*sqrt(t)*Z)
		z := rng.NormFloat64() // Standard normal random variable
		growth := math.Exp((mu-0.5*sigma*sigma)*dt + sigma*math.Sqrt(dt)*z)
		finalValue := config.InitialValue * growth
		results[i] = finalValue

		if finalValue < config.InitialValue {
			lossCount++
		}
	}

	// Sort for percentile calculations
	sort.Float64s(results)

	return MonteCarloResult{
		Mean:               calculateMean(results),
		Median:             percentile(results, 50),
		StandardDeviation:  calculateStdDev(results),
		Percentile5:        percentile(results, 5),
		Percentile25:       percentile(results, 25),
		Percentile75:       percentile(results, 75),
		Percentile95:       percentile(results, 95),
		MinValue:           results[0],
		MaxValue:           results[len(results)-1],
		ProbabilityOfLoss:  float64(lossCount) / float64(config.Simulations) * 100,
		ProbabilityOfGain:  float64(config.Simulations-lossCount) / float64(config.Simulations) * 100,
		Simulations:        config.Simulations,
		Distribution:       sampleDistribution(results, 100),
	}
}

// BettingMonteCarloConfig contains config for betting simulation.
type BettingMonteCarloConfig struct {
	Simulations    int
	InitialBankroll float64
	NumBets        int
	WinProbability float64  // As percentage (0-100)
	AverageOdds    float64  // Decimal odds
	StakePercent   float64  // Percentage of bankroll per bet
	Seed           int64
}

// BettingMonteCarloResult contains betting simulation results.
type BettingMonteCarloResult struct {
	MonteCarloResult
	AvgFinalBankroll   float64
	AvgMaxDrawdown     float64
	RuinProbability    float64 // Probability of going bust
	DoubleProbability  float64 // Probability of doubling bankroll
}

// RunBettingMonteCarlo simulates betting strategies.
func RunBettingMonteCarlo(config BettingMonteCarloConfig) BettingMonteCarloResult {
	if config.Simulations == 0 {
		config.Simulations = 10000
	}
	if config.InitialBankroll == 0 {
		config.InitialBankroll = 1000
	}
	if config.NumBets == 0 {
		config.NumBets = 100
	}
	if config.StakePercent == 0 {
		config.StakePercent = 2 // 2% per bet
	}

	var rng *rand.Rand
	if config.Seed != 0 {
		rng = rand.New(rand.NewSource(config.Seed))
	} else {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	results := make([]float64, config.Simulations)
	drawdowns := make([]float64, config.Simulations)
	ruinCount := 0
	doubleCount := 0

	winProb := config.WinProbability / 100

	for i := 0; i < config.Simulations; i++ {
		bankroll := config.InitialBankroll
		maxBankroll := bankroll
		maxDrawdown := 0.0

		for bet := 0; bet < config.NumBets; bet++ {
			if bankroll <= 0 {
				break
			}

			stake := bankroll * (config.StakePercent / 100)

			if rng.Float64() < winProb {
				// Win
				bankroll += stake * (config.AverageOdds - 1)
			} else {
				// Lose
				bankroll -= stake
			}

			if bankroll > maxBankroll {
				maxBankroll = bankroll
			}

			dd := (maxBankroll - bankroll) / maxBankroll * 100
			if dd > maxDrawdown {
				maxDrawdown = dd
			}
		}

		results[i] = bankroll
		drawdowns[i] = maxDrawdown

		if bankroll <= 0 {
			ruinCount++
		}
		if bankroll >= config.InitialBankroll*2 {
			doubleCount++
		}
	}

	sort.Float64s(results)

	return BettingMonteCarloResult{
		MonteCarloResult: MonteCarloResult{
			Mean:               calculateMean(results),
			Median:             percentile(results, 50),
			StandardDeviation:  calculateStdDev(results),
			Percentile5:        percentile(results, 5),
			Percentile25:       percentile(results, 25),
			Percentile75:       percentile(results, 75),
			Percentile95:       percentile(results, 95),
			MinValue:           results[0],
			MaxValue:           results[len(results)-1],
			ProbabilityOfLoss:  float64(countBelowThreshold(results, config.InitialBankroll)) / float64(config.Simulations) * 100,
			ProbabilityOfGain:  float64(countAboveThreshold(results, config.InitialBankroll)) / float64(config.Simulations) * 100,
			Simulations:        config.Simulations,
			Distribution:       sampleDistribution(results, 100),
		},
		AvgFinalBankroll:  calculateMean(results),
		AvgMaxDrawdown:    calculateMean(drawdowns),
		RuinProbability:   float64(ruinCount) / float64(config.Simulations) * 100,
		DoubleProbability: float64(doubleCount) / float64(config.Simulations) * 100,
	}
}

// SharpeRatio calculates the Sharpe ratio.
// returns: (mean - riskFreeRate) / stdDev
func SharpeRatio(returns []float64, riskFreeRate float64) float64 {
	if len(returns) == 0 {
		return 0
	}
	mean := calculateMean(returns)
	stdDev := calculateStdDev(returns)
	if stdDev == 0 {
		return 0
	}
	return (mean - riskFreeRate) / stdDev
}

// MaxDrawdown calculates the maximum drawdown from a series of values.
func MaxDrawdown(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}

	maxSoFar := values[0]
	maxDrawdown := 0.0

	for _, v := range values {
		if v > maxSoFar {
			maxSoFar = v
		}
		dd := (maxSoFar - v) / maxSoFar * 100
		if dd > maxDrawdown {
			maxDrawdown = dd
		}
	}

	return maxDrawdown
}

// ValueAtRisk calculates VaR at a given confidence level.
// confidence: e.g., 95 for 95% VaR
func ValueAtRisk(returns []float64, confidence float64) float64 {
	if len(returns) == 0 {
		return 0
	}
	sorted := make([]float64, len(returns))
	copy(sorted, returns)
	sort.Float64s(sorted)
	return percentile(sorted, 100-confidence)
}

// Helper functions

func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateStdDev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	mean := calculateMean(values)
	sumSquares := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(values)-1)
	return math.Sqrt(variance)
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 100 {
		return sorted[len(sorted)-1]
	}
	idx := (p / 100) * float64(len(sorted)-1)
	lower := int(idx)
	upper := lower + 1
	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	weight := idx - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

func sampleDistribution(values []float64, sampleSize int) []float64 {
	if len(values) <= sampleSize {
		return values
	}
	sample := make([]float64, sampleSize)
	step := len(values) / sampleSize
	for i := 0; i < sampleSize; i++ {
		sample[i] = values[i*step]
	}
	return sample
}

func countBelowThreshold(values []float64, threshold float64) int {
	count := 0
	for _, v := range values {
		if v < threshold {
			count++
		}
	}
	return count
}

func countAboveThreshold(values []float64, threshold float64) int {
	count := 0
	for _, v := range values {
		if v > threshold {
			count++
		}
	}
	return count
}
