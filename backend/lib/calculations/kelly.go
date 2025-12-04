package calculations

import "math"

// KellyResult contains the result of a Kelly criterion calculation.
type KellyResult struct {
	Stake         float64 // Full Kelly stake
	HalfKelly     float64 // Half Kelly stake
	QuarterKelly  float64 // Quarter Kelly stake
	ExpectedValue float64 // Expected value as percentage
	Edge          float64 // Edge as percentage
}

// CalculateKelly calculates the optimal stake using the Kelly Criterion.
// probability: win probability as percentage (0-100)
// odds: decimal odds (e.g., 2.0 for even money)
// bankroll: total bankroll amount
// fractionalKelly: fraction to use (1.0 = full, 0.5 = half, 0.25 = quarter)
func CalculateKelly(probability, odds, bankroll, fractionalKelly float64) KellyResult {
	p := probability / 100 // Convert to decimal
	q := 1 - p
	b := odds - 1 // Net odds (profit per unit stake)

	// Kelly formula: (bp - q) / b
	kellyFraction := (b*p - q) / b

	if kellyFraction < 0 {
		kellyFraction = 0
	}

	fullKellyStake := kellyFraction * bankroll * fractionalKelly
	edge := (p * odds) - 1
	expectedValue := edge * 100

	return KellyResult{
		Stake:         math.Max(0, fullKellyStake),
		HalfKelly:     math.Max(0, fullKellyStake*0.5),
		QuarterKelly:  math.Max(0, fullKellyStake*0.25),
		ExpectedValue: expectedValue,
		Edge:          edge * 100,
	}
}

// CalculateFullKelly calculates the full Kelly fraction.
// probability: win probability (0 to 1)
// odds: decimal odds
// Returns: Kelly stake as a fraction of bankroll
func CalculateFullKelly(probability, odds float64) float64 {
	if probability <= 0 || probability >= 1 || odds <= 1 {
		return 0
	}
	p := probability
	q := 1 - p
	b := odds - 1
	kelly := (b*p - q) / b
	return math.Max(0, kelly)
}

// CalculateHalfKelly calculates the half Kelly fraction.
func CalculateHalfKelly(probability, odds float64) float64 {
	return CalculateFullKelly(probability, odds) * 0.5
}

// CalculateQuarterKelly calculates the quarter Kelly fraction.
func CalculateQuarterKelly(probability, odds float64) float64 {
	return CalculateFullKelly(probability, odds) * 0.25
}

// CalculateOptimalStake calculates optimal stake with a maximum cap.
func CalculateOptimalStake(probability, odds, bankroll, maxStakePercent float64) float64 {
	kelly := CalculateKelly(probability, odds, bankroll, 1.0)
	maxStake := bankroll * (maxStakePercent / 100)
	return math.Min(kelly.HalfKelly, maxStake)
}

// CalculateImpliedProbability converts decimal odds to implied probability.
// decimalOdds: decimal odds (e.g., 2.00)
// Returns: implied probability as percentage (0-100)
func CalculateImpliedProbability(decimalOdds float64) float64 {
	if decimalOdds <= 1 {
		return 0
	}
	return (1 / decimalOdds) * 100
}

// ProbabilityToOdds converts probability to fair decimal odds.
// probability: probability as percentage (0-100)
// Returns: fair decimal odds
func ProbabilityToOdds(probability float64) float64 {
	if probability <= 0 || probability >= 100 {
		return 0
	}
	return 100 / probability
}

// ValueBetResult contains value bet analysis.
type ValueBetResult struct {
	ImpliedProbability float64
	Value              float64
	IsValueBet         bool
	IsHighValue        bool
	ExpectedValue      float64
	Recommendation     string // "skip", "bet", "strong_bet"
}

// DetectValueBet analyzes if a bet has value.
// trueProbability: true probability as percentage (0-100)
// bookmakerOdds: decimal odds from bookmaker
// valueThreshold: minimum value percentage to consider (default 5%)
func DetectValueBet(trueProbability, bookmakerOdds, valueThreshold float64) ValueBetResult {
	if valueThreshold == 0 {
		valueThreshold = 5
	}

	impliedProbability := CalculateImpliedProbability(bookmakerOdds)
	value := trueProbability - impliedProbability
	expectedValue := (trueProbability/100)*(bookmakerOdds-1) - (1 - trueProbability/100)

	recommendation := "skip"
	if value > 10 {
		recommendation = "strong_bet"
	} else if value > valueThreshold {
		recommendation = "bet"
	}

	return ValueBetResult{
		ImpliedProbability: impliedProbability,
		Value:              value,
		IsValueBet:         value > valueThreshold,
		IsHighValue:        value > 10,
		ExpectedValue:      expectedValue * 100,
		Recommendation:     recommendation,
	}
}

// CalculateValue computes value percentage between fair odds and bookmaker odds.
func CalculateValue(fairProbability, bookmakerOdds float64) (value float64, isValueBet bool, expectedValue float64) {
	fairOdds := 100 / fairProbability
	value = ((bookmakerOdds / fairOdds) - 1) * 100
	expectedValue = (fairProbability/100)*(bookmakerOdds-1) - (1 - fairProbability/100)
	expectedValue *= 100
	isValueBet = value > 0
	return
}

// BayesianUpdate updates probability based on new evidence.
// priorProbability: prior probability (0-1)
// likelihood: probability of evidence given hypothesis (0-1)
// evidenceProbability: overall probability of evidence (0-1)
// Returns: posterior probability
func BayesianUpdate(priorProbability, likelihood, evidenceProbability float64) float64 {
	if evidenceProbability == 0 {
		return priorProbability
	}
	return (likelihood * priorProbability) / evidenceProbability
}

// CalculateWeightedProbability computes weighted average from multiple models.
// probabilities: array of probability estimates (0-100)
// weights: array of weights for each model (should sum to 1)
func CalculateWeightedProbability(probabilities, weights []float64) (float64, error) {
	if len(probabilities) != len(weights) {
		return 0, nil // Return 0 if lengths don't match
	}

	totalWeight := 0.0
	for _, w := range weights {
		totalWeight += w
	}

	// Normalize weights if they don't sum to 1
	normalizedWeights := make([]float64, len(weights))
	for i, w := range weights {
		normalizedWeights[i] = w / totalWeight
	}

	result := 0.0
	for i, prob := range probabilities {
		result += prob * normalizedWeights[i]
	}
	return result, nil
}

// CalculateEnsembleProbability combines probabilities from multiple models.
func CalculateEnsembleProbability(poissonProb, eloProb, statProb float64, xgProb *float64) float64 {
	probs := []float64{poissonProb, eloProb, statProb}
	if xgProb != nil {
		probs = append(probs, *xgProb)
	}
	equalWeight := 1.0 / float64(len(probs))
	weights := make([]float64, len(probs))
	for i := range weights {
		weights[i] = equalWeight
	}
	result, _ := CalculateWeightedProbability(probs, weights)
	return result
}

// ArbitrageResult contains arbitrage analysis.
type ArbitrageResult struct {
	IsArbitrage bool
	Margin      float64
	BestHome    OddsSelection
	BestDraw    OddsSelection
	BestAway    OddsSelection
	Stakes      *ArbitrageStakes
}

// OddsSelection represents the best odds selection.
type OddsSelection struct {
	Index int
	Odds  float64
}

// ArbitrageStakes contains stake distribution for arbitrage.
type ArbitrageStakes struct {
	Home float64
	Draw float64
	Away float64
}

// FindArbitrage detects arbitrage opportunities across bookmakers.
func FindArbitrage(homeOdds, drawOdds, awayOdds []float64) ArbitrageResult {
	bestHome := OddsSelection{Index: 0, Odds: homeOdds[0]}
	for i, o := range homeOdds {
		if o > bestHome.Odds {
			bestHome = OddsSelection{Index: i, Odds: o}
		}
	}

	bestDraw := OddsSelection{Index: 0, Odds: drawOdds[0]}
	for i, o := range drawOdds {
		if o > bestDraw.Odds {
			bestDraw = OddsSelection{Index: i, Odds: o}
		}
	}

	bestAway := OddsSelection{Index: 0, Odds: awayOdds[0]}
	for i, o := range awayOdds {
		if o > bestAway.Odds {
			bestAway = OddsSelection{Index: i, Odds: o}
		}
	}

	totalImplied := (1 / bestHome.Odds) + (1 / bestDraw.Odds) + (1 / bestAway.Odds)
	margin := (1 - totalImplied) * 100
	isArbitrage := totalImplied < 1

	var stakes *ArbitrageStakes
	if isArbitrage {
		total := 100.0
		stakes = &ArbitrageStakes{
			Home: (total / bestHome.Odds) / totalImplied,
			Draw: (total / bestDraw.Odds) / totalImplied,
			Away: (total / bestAway.Odds) / totalImplied,
		}
	}

	return ArbitrageResult{
		IsArbitrage: isArbitrage,
		Margin:      margin,
		BestHome:    bestHome,
		BestDraw:    bestDraw,
		BestAway:    bestAway,
		Stakes:      stakes,
	}
}

// KellyGrowthResult contains simulation results.
type KellyGrowthResult struct {
	FinalBankroll float64
	Growth        float64
	MaxDrawdown   float64
}

// BetOutcome represents a bet with its outcome.
type BetOutcome struct {
	Probability float64
	Odds        float64
	Won         bool
}

// SimulateKellyGrowth simulates bankroll growth using Kelly staking.
func SimulateKellyGrowth(initialBankroll float64, bets []BetOutcome, fraction float64) KellyGrowthResult {
	if fraction == 0 {
		fraction = 0.5 // Default to half Kelly
	}

	bankroll := initialBankroll
	maxBankroll := initialBankroll
	maxDrawdown := 0.0

	for _, bet := range bets {
		kelly := CalculateKelly(bet.Probability, bet.Odds, bankroll, fraction)
		stake := kelly.Stake

		if bet.Won {
			bankroll += stake * (bet.Odds - 1)
		} else {
			bankroll -= stake
		}

		if bankroll > maxBankroll {
			maxBankroll = bankroll
		}

		drawdown := ((maxBankroll - bankroll) / maxBankroll) * 100
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return KellyGrowthResult{
		FinalBankroll: bankroll,
		Growth:        ((bankroll - initialBankroll) / initialBankroll) * 100,
		MaxDrawdown:   maxDrawdown,
	}
}
