package calculations

import (
	"math"
	"math/rand"
)

// NormalDistribution calculates PDF of normal distribution at x.
func NormalDistribution(x, mean, stdDev float64) float64 {
	coefficient := 1 / (stdDev * math.Sqrt(2*math.Pi))
	exponent := -math.Pow(x-mean, 2) / (2 * math.Pow(stdDev, 2))
	return coefficient * math.Exp(exponent)
}

// NormalCDF calculates cumulative distribution function for normal distribution.
// Approximation using error function
func NormalCDF(x, mean, stdDev float64) float64 {
	z := (x - mean) / (stdDev * math.Sqrt(2))
	return 0.5 * (1 + math.Erf(z))
}

// ZScore calculates z-score (standard score).
func ZScore(x, mean, stdDev float64) float64 {
	if stdDev == 0 {
		return 0
	}
	return (x - mean) / stdDev
}

// ConfidenceIntervalNormal calculates confidence interval for normal distribution.
func ConfidenceIntervalNormal(mean, stdDev, confidenceLevel float64, sampleSize int) (float64, float64) {
	var zScore float64
	switch confidenceLevel {
	case 0.90:
		zScore = 1.645
	case 0.95:
		zScore = 1.960
	case 0.99:
		zScore = 2.576
	default:
		zScore = 1.960 // Default to 95%
	}

	marginOfError := zScore * (stdDev / math.Sqrt(float64(sampleSize)))
	return mean - marginOfError, mean + marginOfError
}

// BayesianUpdate updates probability using Bayes' theorem.
// P(A|B) = P(B|A) * P(A) / P(B)
func BayesianUpdate(priorProbability, likelihood, evidence float64) float64 {
	if evidence == 0 {
		return 0
	}
	return (likelihood * priorProbability) / evidence
}

// BayesianEvidence calculates total probability P(B).
// P(B) = P(B|A) * P(A) + P(B|not A) * P(not A)
func BayesianEvidence(likelihoodA, priorA, likelihoodNotA float64) float64 {
	return (likelihoodA * priorA) + (likelihoodNotA * (1 - priorA))
}

// BinomialProbability calculates probability of k successes in n trials.
func BinomialProbability(n, k int, p float64) float64 {
	if k > n || k < 0 {
		return 0
	}

	// Calculate binomial coefficient C(n,k)
	coefficient := float64(binomialCoefficient(n, k))

	// Calculate probability
	return coefficient * math.Pow(p, float64(k)) * math.Pow(1-p, float64(n-k))
}

func binomialCoefficient(n, k int) int {
	if k > n-k {
		k = n - k
	}

	result := 1
	for i := 0; i < k; i++ {
		result *= (n - i)
		result /= (i + 1)
	}

	return result
}

// BinomialMean calculates expected value of binomial distribution.
func BinomialMean(n int, p float64) float64 {
	return float64(n) * p
}

// BinomialVariance calculates variance of binomial distribution.
func BinomialVariance(n int, p float64) float64 {
	return float64(n) * p * (1 - p)
}

// GeometricProbability calculates probability of first success on trial k.
func GeometricProbability(k int, p float64) float64 {
	if k < 1 {
		return 0
	}
	return math.Pow(1-p, float64(k-1)) * p
}

// GeometricMean calculates expected number of trials until first success.
func GeometricMean(p float64) float64 {
	if p == 0 {
		return 0
	}
	return 1 / p
}

// HypergeometricProbability calculates probability without replacement.
// P(X = k) where drawing k successes from population
func HypergeometricProbability(populationSize, successInPopulation, sampleSize, successInSample int) float64 {
	numerator := float64(binomialCoefficient(successInPopulation, successInSample) *
		binomialCoefficient(populationSize-successInPopulation, sampleSize-successInSample))

	denominator := float64(binomialCoefficient(populationSize, sampleSize))

	if denominator == 0 {
		return 0
	}

	return numerator / denominator
}

// ExponentialDistribution calculates PDF of exponential distribution.
func ExponentialDistribution(x, lambda float64) float64 {
	if x < 0 || lambda <= 0 {
		return 0
	}
	return lambda * math.Exp(-lambda*x)
}

// ExponentialCDF calculates CDF of exponential distribution.
func ExponentialCDF(x, lambda float64) float64 {
	if x < 0 {
		return 0
	}
	return 1 - math.Exp(-lambda*x)
}

// LogNormalDistribution calculates PDF of log-normal distribution.
func LogNormalDistribution(x, mu, sigma float64) float64 {
	if x <= 0 || sigma <= 0 {
		return 0
	}

	coefficient := 1 / (x * sigma * math.Sqrt(2*math.Pi))
	exponent := -math.Pow(math.Log(x)-mu, 2) / (2 * math.Pow(sigma, 2))

	return coefficient * math.Exp(exponent)
}

// StudentTDistribution approximates t-distribution PDF.
func StudentTDistribution(x float64, degreesOfFreedom int) float64 {
	v := float64(degreesOfFreedom)

	coefficient := math.Gamma((v+1)/2) / (math.Sqrt(v*math.Pi) * math.Gamma(v/2))
	power := -((v + 1) / 2)

	return coefficient * math.Pow(1+(x*x)/v, power)
}

// ChiSquareDistribution calculates PDF of chi-square distribution.
func ChiSquareDistribution(x float64, k int) float64 {
	if x < 0 || k < 1 {
		return 0
	}

	kHalf := float64(k) / 2
	coefficient := math.Pow(x, kHalf-1) * math.Exp(-x/2)
	denominator := math.Pow(2, kHalf) * math.Gamma(kHalf)

	if denominator == 0 {
		return 0
	}

	return coefficient / denominator
}

// MonteCarloSimulation runs Monte Carlo simulation for price prediction.
type MonteCarloParams struct {
	InitialPrice  float64
	DriftRate     float64 // Expected return
	Volatility    float64 // Standard deviation
	DaysToSimulate int
	NumSimulations int
}

type MonteCarloResult struct {
	Mean          float64
	Median        float64
	StdDev        float64
	Percentile5   float64
	Percentile95  float64
	ProbAbove     float64 // Probability of price being above initial
	AllPrices     []float64
}

func MonteCarloSimulation(params MonteCarloParams) MonteCarloResult {
	finalPrices := make([]float64, params.NumSimulations)

	for i := 0; i < params.NumSimulations; i++ {
		price := params.InitialPrice
		dt := 1.0 // Daily time step

		for day := 0; day < params.DaysToSimulate; day++ {
			// Generate random normal variable
			z := rand.NormFloat64()

			// Geometric Brownian Motion
			drift := (params.DriftRate - 0.5*params.Volatility*params.Volatility) * dt
			diffusion := params.Volatility * math.Sqrt(dt) * z

			price *= math.Exp(drift + diffusion)
		}

		finalPrices[i] = price
	}

	// Calculate statistics
	result := MonteCarloResult{AllPrices: finalPrices}

	sum := 0.0
	aboveCount := 0

	for _, price := range finalPrices {
		sum += price
		if price > params.InitialPrice {
			aboveCount++
		}
	}

	result.Mean = sum / float64(params.NumSimulations)
	result.ProbAbove = float64(aboveCount) / float64(params.NumSimulations)

	// Calculate standard deviation
	variance := 0.0
	for _, price := range finalPrices {
		variance += math.Pow(price-result.Mean, 2)
	}
	result.StdDev = math.Sqrt(variance / float64(params.NumSimulations))

	// For percentiles, we'd need to sort (simplified here)
	result.Percentile5 = result.Mean - 1.645*result.StdDev
	result.Percentile95 = result.Mean + 1.645*result.StdDev
	result.Median = result.Mean // Approximation

	return result
}

// LinearRegression performs simple linear regression.
type RegressionResult struct {
	Slope       float64
	Intercept   float64
	RSquared    float64
	Correlation float64
}

func LinearRegression(x, y []float64) RegressionResult {
	if len(x) != len(y) || len(x) == 0 {
		return RegressionResult{}
	}

	n := float64(len(x))

	// Calculate means
	sumX, sumY := 0.0, 0.0
	for i := range x {
		sumX += x[i]
		sumY += y[i]
	}
	meanX := sumX / n
	meanY := sumY / n

	// Calculate slope and intercept
	numerator := 0.0
	denominator := 0.0

	for i := range x {
		numerator += (x[i] - meanX) * (y[i] - meanY)
		denominator += (x[i] - meanX) * (x[i] - meanX)
	}

	slope := 0.0
	if denominator != 0 {
		slope = numerator / denominator
	}

	intercept := meanY - slope*meanX

	// Calculate R-squared
	ssTotal := 0.0
	ssResidual := 0.0

	for i := range x {
		predicted := slope*x[i] + intercept
		ssTotal += math.Pow(y[i]-meanY, 2)
		ssResidual += math.Pow(y[i]-predicted, 2)
	}

	rSquared := 0.0
	if ssTotal != 0 {
		rSquared = 1 - (ssResidual / ssTotal)
	}

	// Calculate correlation coefficient
	correlation := math.Sqrt(rSquared)
	if slope < 0 {
		correlation = -correlation
	}

	return RegressionResult{
		Slope:       slope,
		Intercept:   intercept,
		RSquared:    rSquared,
		Correlation: correlation,
	}
}

// TTest performs one-sample t-test.
type TTestResult struct {
	TStatistic float64
	PValue     float64 // Approximate
	Significant bool
}

func TTest(sample []float64, hypothesizedMean, alpha float64) TTestResult {
	if len(sample) == 0 {
		return TTestResult{}
	}

	n := float64(len(sample))

	// Calculate sample mean
	sum := 0.0
	for _, value := range sample {
		sum += value
	}
	sampleMean := sum / n

	// Calculate sample standard deviation
	variance := 0.0
	for _, value := range sample {
		variance += math.Pow(value-sampleMean, 2)
	}
	stdDev := math.Sqrt(variance / (n - 1))

	// Calculate t-statistic
	standardError := stdDev / math.Sqrt(n)
	tStatistic := 0.0
	if standardError != 0 {
		tStatistic = (sampleMean - hypothesizedMean) / standardError
	}

	// Approximate p-value (two-tailed)
	// This is simplified - in production, use proper t-distribution
	pValue := 2 * (1 - NormalCDF(math.Abs(tStatistic), 0, 1))

	return TTestResult{
		TStatistic:  tStatistic,
		PValue:      pValue,
		Significant: pValue < alpha,
	}
}

// MovingAverageConvergence calculates probability of trend continuation.
func MovingAverageConvergence(shortMA, longMA, currentPrice float64) float64 {
	if longMA == 0 {
		return 0.5
	}

	// Distance from moving averages
	shortDistance := (currentPrice - shortMA) / shortMA
	longDistance := (currentPrice - longMA) / longMA

	// Convergence indicator
	convergence := (shortMA - longMA) / longMA

	// Combine signals (simplified probability model)
	signal := (shortDistance + longDistance + convergence) / 3

	// Convert to probability (0-1 range)
	probability := 0.5 + (signal * 0.5)

	// Clamp to valid probability range
	if probability < 0 {
		probability = 0
	}
	if probability > 1 {
		probability = 1
	}

	return probability
}

// KernelDensityEstimation estimates probability density.
func KernelDensityEstimation(data []float64, point, bandwidth float64) float64 {
	if len(data) == 0 || bandwidth == 0 {
		return 0
	}

	n := float64(len(data))
	sum := 0.0

	for _, x := range data {
		u := (point - x) / bandwidth
		// Gaussian kernel
		kernel := math.Exp(-0.5*u*u) / math.Sqrt(2*math.Pi)
		sum += kernel
	}

	return sum / (n * bandwidth)
}

// BootstrapConfidenceInterval calculates CI using bootstrap method.
func BootstrapConfidenceInterval(data []float64, numBootstrap int, confidenceLevel float64) (float64, float64) {
	if len(data) == 0 {
		return 0, 0
	}

	means := make([]float64, numBootstrap)

	for i := 0; i < numBootstrap; i++ {
		// Resample with replacement
		sample := make([]float64, len(data))
		for j := range sample {
			sample[j] = data[rand.Intn(len(data))]
		}

		// Calculate mean of bootstrap sample
		sum := 0.0
		for _, value := range sample {
			sum += value
		}
		means[i] = sum / float64(len(sample))
	}

	// Sort means (simplified - in production use proper sorting)
	// Return approximate percentiles
	alpha := 1 - confidenceLevel
	lowerPercentile := alpha / 2
	upperPercentile := 1 - (alpha / 2)

	// Simplified calculation
	sum := 0.0
	for _, mean := range means {
		sum += mean
	}
	avgMean := sum / float64(numBootstrap)

	// Use standard deviation for interval
	variance := 0.0
	for _, mean := range means {
		variance += math.Pow(mean-avgMean, 2)
	}
	stdDev := math.Sqrt(variance / float64(numBootstrap))

	return avgMean - 1.96*stdDev, avgMean + 1.96*stdDev
}

// ProbabilityOfProfit calculates probability of profit for a trade.
func ProbabilityOfProfit(entryPrice, targetPrice, stdDev float64, daysToTarget int) float64 {
	if stdDev == 0 || daysToTarget == 0 {
		return 0.5
	}

	// Calculate required return
	requiredReturn := (targetPrice - entryPrice) / entryPrice

	// Annualize volatility
	annualVolatility := stdDev * math.Sqrt(252.0/float64(daysToTarget))

	// Calculate z-score
	z := requiredReturn / annualVolatility

	// Return probability using normal CDF
	return NormalCDF(z, 0, 1)
}
