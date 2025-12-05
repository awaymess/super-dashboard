package calculations

import (
	"math"
)

// KellyCriterion calculates the optimal bet size using Kelly Criterion.
// f* = (bp - q) / b
// where b = odds - 1, p = win probability, q = 1 - p
func KellyCriterion(bankroll, odds, winProbability float64) float64 {
	b := odds - 1
	p := winProbability
	q := 1 - p

	kellyFraction := (b*p - q) / b

	// Don't bet if Kelly is negative (no edge)
	if kellyFraction <= 0 {
		return 0
	}

	return bankroll * kellyFraction
}

// FractionalKelly applies a fraction of Kelly for safety.
// Common fractions: 0.25 (quarter), 0.5 (half)
func FractionalKelly(bankroll, odds, winProbability, fraction float64) float64 {
	fullKelly := KellyCriterion(bankroll, odds, winProbability)
	return fullKelly * fraction
}

// ExpectedValue calculates the expected value of a bet.
// EV = (probability of winning × amount won per bet) - (probability of losing × amount lost per bet)
func ExpectedValue(stake, odds, winProbability float64) float64 {
	potentialWin := stake * (odds - 1)
	potentialLoss := stake
	loseProbability := 1 - winProbability

	ev := (winProbability * potentialWin) - (loseProbability * potentialLoss)
	return ev
}

// ExpectedValuePercent calculates EV as a percentage of stake.
func ExpectedValuePercent(stake, odds, winProbability float64) float64 {
	ev := ExpectedValue(stake, odds, winProbability)
	return (ev / stake) * 100
}

// ImpliedProbability calculates the implied probability from odds.
func ImpliedProbability(odds float64) float64 {
	return 1 / odds
}

// TrueProbabilityFromOdds removes bookmaker margin to get true probability.
func TrueProbabilityFromOdds(homeOdds, drawOdds, awayOdds float64) (float64, float64, float64) {
	totalImplied := ImpliedProbability(homeOdds) + ImpliedProbability(drawOdds) + ImpliedProbability(awayOdds)
	margin := totalImplied - 1

	// Remove margin proportionally
	homeProb := ImpliedProbability(homeOdds) / totalImplied
	drawProb := ImpliedProbability(drawOdds) / totalImplied
	awayProb := ImpliedProbability(awayOdds) / totalImplied

	return homeProb, drawProb, awayProb
}

// ValueBetPercentage calculates value bet percentage.
// Value = (True Probability × Odds) - 1
func ValueBetPercentage(trueProbability, odds float64) float64 {
	value := (trueProbability * odds) - 1
	return value * 100
}

// PoissonProbability calculates Poisson probability for exact k events.
// P(k) = (λ^k × e^(-λ)) / k!
func PoissonProbability(lambda float64, k int) float64 {
	if lambda <= 0 {
		return 0
	}
	return (math.Pow(lambda, float64(k)) * math.Exp(-lambda)) / factorial(float64(k))
}

// PoissonUnderGoals calculates probability of under X goals.
func PoissonUnderGoals(expectedGoals float64, threshold int) float64 {
	probability := 0.0
	for k := 0; k <= threshold; k++ {
		probability += PoissonProbability(expectedGoals, k)
	}
	return probability
}

// PoissonOverGoals calculates probability of over X goals.
func PoissonOverGoals(expectedGoals float64, threshold int) float64 {
	return 1 - PoissonUnderGoals(expectedGoals, threshold)
}

// PoissonExactScore calculates probability of exact score.
func PoissonExactScore(homeExpected, awayExpected float64, homeGoals, awayGoals int) float64 {
	homeProb := PoissonProbability(homeExpected, homeGoals)
	awayProb := PoissonProbability(awayExpected, awayGoals)
	return homeProb * awayProb
}

// PoissonCorrectScore calculates probabilities for common scores.
func PoissonCorrectScore(homeExpected, awayExpected float64) map[string]float64 {
	scores := make(map[string]float64)
	
	// Calculate probabilities for scores 0-0 to 5-5
	for h := 0; h <= 5; h++ {
		for a := 0; a <= 5; a++ {
			scoreKey := string(rune('0'+h)) + "-" + string(rune('0'+a))
			scores[scoreKey] = PoissonExactScore(homeExpected, awayExpected, h, a)
		}
	}
	
	return scores
}

// ELOExpectedScore calculates expected score based on ELO ratings.
// E = 1 / (1 + 10^((R_opponent - R_player) / 400))
func ELOExpectedScore(playerRating, opponentRating float64) float64 {
	return 1 / (1 + math.Pow(10, (opponentRating-playerRating)/400))
}

// ELOWinProbability calculates win probability from ELO ratings.
func ELOWinProbability(rating1, rating2 float64) float64 {
	return ELOExpectedScore(rating1, rating2)
}

// ELODrawProbability estimates draw probability from ELO difference.
func ELODrawProbability(rating1, rating2 float64) float64 {
	diff := math.Abs(rating1 - rating2)
	// Empirical formula: closer ratings = higher draw probability
	baseDraw := 0.25
	adjustment := (200 - diff) / 1000
	if adjustment < 0 {
		adjustment = 0
	}
	return baseDraw + adjustment
}

// ELO1X2Probabilities calculates home/draw/away probabilities from ELO.
func ELO1X2Probabilities(homeRating, awayRating float64) (float64, float64, float64) {
	homeWin := ELOWinProbability(homeRating, awayRating)
	drawProb := ELODrawProbability(homeRating, awayRating)
	awayWin := 1 - homeWin - drawProb
	
	if awayWin < 0 {
		awayWin = 0
	}
	
	// Normalize to sum to 1
	total := homeWin + drawProb + awayWin
	return homeWin / total, drawProb / total, awayWin / total
}

// ClosingLineValue calculates CLV (Closing Line Value).
// CLV = (Closing Odds / Opening Odds) - 1
func ClosingLineValue(openingOdds, closingOdds float64) float64 {
	return ((closingOdds / openingOdds) - 1) * 100
}

// BreakEvenPoint calculates the win rate needed to break even.
func BreakEvenPoint(odds float64) float64 {
	return (1 / odds) * 100
}

// ROI calculates Return on Investment.
func ROI(totalStake, totalProfit float64) float64 {
	if totalStake == 0 {
		return 0
	}
	return (totalProfit / totalStake) * 100
}

// Yield calculates betting yield (average ROI per bet).
func Yield(totalStake, totalProfit float64, numberOfBets int) float64 {
	if numberOfBets == 0 {
		return 0
	}
	roi := ROI(totalStake, totalProfit)
	return roi / float64(numberOfBets)
}

// AverageOdds calculates the average odds from a slice of odds.
func AverageOdds(odds []float64) float64 {
	if len(odds) == 0 {
		return 0
	}
	sum := 0.0
	for _, o := range odds {
		sum += o
	}
	return sum / float64(len(odds))
}

// BookmakerMargin calculates the overround/margin in odds.
func BookmakerMargin(odds []float64) float64 {
	totalImplied := 0.0
	for _, o := range odds {
		totalImplied += ImpliedProbability(o)
	}
	return (totalImplied - 1) * 100
}

// FairOdds removes bookmaker margin to get fair odds.
func FairOdds(odds, totalImplied float64) float64 {
	impliedProb := ImpliedProbability(odds)
	trueProb := impliedProb / totalImplied
	return 1 / trueProb
}

// ArbitrageProfit calculates profit from arbitrage betting.
func ArbitrageProfit(stake float64, odds []float64) float64 {
	totalImplied := 0.0
	for _, o := range odds {
		totalImplied += ImpliedProbability(o)
	}
	
	// No arbitrage if total implied >= 1
	if totalImplied >= 1 {
		return 0
	}
	
	// Calculate profit
	return (1/totalImplied - 1) * stake
}

// ArbitrageStakes calculates individual stakes for arbitrage.
func ArbitrageStakes(totalStake float64, odds []float64) []float64 {
	stakes := make([]float64, len(odds))
	totalImplied := 0.0
	
	for _, o := range odds {
		totalImplied += ImpliedProbability(o)
	}
	
	// Calculate stakes proportionally
	for i, o := range odds {
		stakes[i] = (totalStake * ImpliedProbability(o)) / totalImplied
	}
	
	return stakes
}

// CompoundGrowth calculates compound growth over time.
func CompoundGrowth(initialBankroll, roi float64, numberOfBets int) float64 {
	rate := 1 + (roi / 100)
	return initialBankroll * math.Pow(rate, float64(numberOfBets))
}

// VarianceCalculation calculates variance for a series of bets.
func VarianceCalculation(results []float64) float64 {
	if len(results) == 0 {
		return 0
	}
	
	mean := 0.0
	for _, r := range results {
		mean += r
	}
	mean /= float64(len(results))
	
	variance := 0.0
	for _, r := range results {
		variance += math.Pow(r-mean, 2)
	}
	variance /= float64(len(results))
	
	return variance
}

// StandardDeviation calculates standard deviation.
func StandardDeviation(results []float64) float64 {
	return math.Sqrt(VarianceCalculation(results))
}

// ConfidenceInterval calculates 95% confidence interval.
func ConfidenceInterval(mean, stdDev float64, sampleSize int) (float64, float64) {
	// 1.96 for 95% confidence
	margin := 1.96 * (stdDev / math.Sqrt(float64(sampleSize)))
	return mean - margin, mean + margin
}

// factorial calculates factorial (helper function).
func factorial(n float64) float64 {
	if n <= 1 {
		return 1
	}
	result := 1.0
	for i := 2.0; i <= n; i++ {
		result *= i
	}
	return result
}

// BettingBankrollGrowth models bankroll growth with Kelly betting.
func BettingBankrollGrowth(initialBankroll, winRate, avgOdds float64, numberOfBets int, fraction float64) float64 {
	bankroll := initialBankroll
	
	for i := 0; i < numberOfBets; i++ {
		stake := FractionalKelly(bankroll, avgOdds, winRate, fraction)
		
		// Simulate win/loss (simplified)
		if winRate >= 0.5 { // Win
			bankroll += stake * (avgOdds - 1)
		} else { // Loss
			bankroll -= stake
		}
	}
	
	return bankroll
}

// OptimalKellyFraction finds optimal Kelly fraction through simulation.
func OptimalKellyFraction(bankroll, winRate, avgOdds float64, numberOfTrials int) float64 {
	fractions := []float64{0.1, 0.25, 0.5, 0.75, 1.0}
	bestFraction := 0.25
	bestResult := 0.0
	
	for _, fraction := range fractions {
		avgResult := 0.0
		for i := 0; i < numberOfTrials; i++ {
			result := BettingBankrollGrowth(bankroll, winRate, avgOdds, 100, fraction)
			avgResult += result
		}
		avgResult /= float64(numberOfTrials)
		
		if avgResult > bestResult {
			bestResult = avgResult
			bestFraction = fraction
		}
	}
	
	return bestFraction
}
