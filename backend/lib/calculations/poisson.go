// Package calculations provides mathematical models for betting analytics and stock valuation.
// These functions mirror the frontend TypeScript implementations in frontend/lib/calculations.
package calculations

import "math"

// Factorial calculates n! for Poisson calculations.
func Factorial(n int) float64 {
	if n <= 1 {
		return 1
	}
	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	return result
}

// PoissonProbability calculates the Poisson probability P(X=k) for a given lambda.
// P(X=k) = (λ^k * e^-λ) / k!
func PoissonProbability(lambda float64, k int) float64 {
	if k < 0 {
		return 0
	}
	if lambda == 0 {
		if k == 0 {
			return 1 // P(X=0|λ=0) = 1
		}
		return 0
	}
	if lambda < 0 {
		return 0
	}
	return math.Pow(lambda, float64(k)) * math.Exp(-lambda) / Factorial(k)
}

// PoissonPrediction contains the result of a Poisson-based match prediction.
type PoissonPrediction struct {
	ExpectedHomeGoals float64
	ExpectedAwayGoals float64
	HomeWinProb       float64
	DrawProb          float64
	AwayWinProb       float64
	Over25Prob        float64
	Under25Prob       float64
	BTTSProb          float64
	MostLikelyScores  []ScoreProbability
}

// ScoreProbability represents a score and its probability.
type ScoreProbability struct {
	HomeGoals   int
	AwayGoals   int
	Probability float64
}

// CalculatePoissonPrediction computes match outcome probabilities using the Poisson distribution.
// It takes team attack/defense strengths and league average goals to predict outcomes.
func CalculatePoissonPrediction(
	homeGoalsAvg, homeConcededAvg,
	awayGoalsAvg, awayConcededAvg,
	leagueAvgGoals float64,
) PoissonPrediction {
	if leagueAvgGoals == 0 {
		leagueAvgGoals = 2.75 // Default league average
	}

	halfLeagueAvg := leagueAvgGoals / 2

	// Calculate attack and defense strengths
	homeAttackStrength := homeGoalsAvg / halfLeagueAvg
	homeDefenseStrength := homeConcededAvg / halfLeagueAvg
	awayAttackStrength := awayGoalsAvg / halfLeagueAvg
	awayDefenseStrength := awayConcededAvg / halfLeagueAvg

	// Calculate expected goals
	expectedHomeGoals := homeAttackStrength * awayDefenseStrength * halfLeagueAvg
	expectedAwayGoals := awayAttackStrength * homeDefenseStrength * halfLeagueAvg

	// Calculate probabilities using score matrix
	maxGoals := 10
	var homeWinProb, drawProb, awayWinProb float64
	var over25Prob, bttsProb float64

	scores := make([]ScoreProbability, 0)

	for i := 0; i <= maxGoals; i++ {
		for j := 0; j <= maxGoals; j++ {
			prob := PoissonProbability(expectedHomeGoals, i) * PoissonProbability(expectedAwayGoals, j)

			if i > j {
				homeWinProb += prob
			} else if i == j {
				drawProb += prob
			} else {
				awayWinProb += prob
			}

			if i+j > 2 {
				over25Prob += prob
			}
			if i > 0 && j > 0 {
				bttsProb += prob
			}

			if i <= 5 && j <= 5 {
				scores = append(scores, ScoreProbability{
					HomeGoals:   i,
					AwayGoals:   j,
					Probability: prob * 100,
				})
			}
		}
	}

	// Sort scores by probability (descending)
	sortScoresByProbability(scores)

	// Take top 10
	if len(scores) > 10 {
		scores = scores[:10]
	}

	return PoissonPrediction{
		ExpectedHomeGoals: expectedHomeGoals,
		ExpectedAwayGoals: expectedAwayGoals,
		HomeWinProb:       homeWinProb * 100,
		DrawProb:          drawProb * 100,
		AwayWinProb:       awayWinProb * 100,
		Over25Prob:        over25Prob * 100,
		Under25Prob:       (1 - over25Prob) * 100,
		BTTSProb:          bttsProb * 100,
		MostLikelyScores:  scores,
	}
}

// sortScoresByProbability sorts scores by probability in descending order.
func sortScoresByProbability(scores []ScoreProbability) {
	// Simple bubble sort for small arrays
	n := len(scores)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if scores[j].Probability < scores[j+1].Probability {
				scores[j], scores[j+1] = scores[j+1], scores[j]
			}
		}
	}
}

// CalculateGoalProbabilities computes outcome probabilities from expected goals.
func CalculateGoalProbabilities(homeExpectedGoals, awayExpectedGoals float64) (homeWin, draw, awayWin, over25, under25, btts float64) {
	maxGoals := 10

	for i := 0; i <= maxGoals; i++ {
		for j := 0; j <= maxGoals; j++ {
			prob := PoissonProbability(homeExpectedGoals, i) * PoissonProbability(awayExpectedGoals, j)
			if i > j {
				homeWin += prob
			} else if i == j {
				draw += prob
			} else {
				awayWin += prob
			}
			if i+j > 2 {
				over25 += prob
			}
			if i > 0 && j > 0 {
				btts += prob
			}
		}
	}
	under25 = 1 - over25
	return
}

// CalculatePoissonProbabilities returns probability distribution for 0 to maxGoals.
func CalculatePoissonProbabilities(lambda float64, maxGoals int) []float64 {
	probs := make([]float64, maxGoals+1)
	for k := 0; k <= maxGoals; k++ {
		probs[k] = PoissonProbability(lambda, k)
	}
	return probs
}
