package calculations

import "math"

const (
	// BaseELO is the starting ELO rating for new teams.
	BaseELO = 1500
	// KFactor is the default K-factor for ELO updates.
	KFactor = 32
	// HomeAdvantage is the default home advantage in ELO points.
	HomeAdvantage = 100
)

// ELORating represents a team's ELO rating and recent change.
type ELORating struct {
	Rating int
	Change int
}

// MatchResult represents a match outcome for ELO calculations.
type MatchResult struct {
	HomeTeam  string
	AwayTeam  string
	HomeScore int
	AwayScore int
}

// ELOMatchProbabilities represents win/draw/loss probabilities.
type ELOMatchProbabilities struct {
	HomeWin float64
	Draw    float64
	AwayWin float64
}

// CalculateExpectedScore calculates the expected score based on ELO ratings.
// This is the probability that player A beats player B.
func CalculateExpectedScore(ratingA, ratingB float64) float64 {
	return 1 / (1 + math.Pow(10, (ratingB-ratingA)/400))
}

// CalculateNewRating computes the new ELO rating after a match.
// rating: current rating
// expectedScore: expected outcome (0-1)
// actualScore: actual outcome (1=win, 0.5=draw, 0=loss)
// kFactor: K-factor for sensitivity (default 32)
func CalculateNewRating(rating, expectedScore, actualScore float64, kFactor float64) float64 {
	if kFactor == 0 {
		kFactor = KFactor
	}
	return rating + kFactor*(actualScore-expectedScore)
}

// CalculateELOMatchProbabilities predicts match outcome probabilities from ELO ratings.
func CalculateELOMatchProbabilities(homeRating, awayRating float64, homeAdvantage float64) ELOMatchProbabilities {
	if homeAdvantage == 0 {
		homeAdvantage = HomeAdvantage
	}

	adjustedHomeRating := homeRating + homeAdvantage
	expectedHome := CalculateExpectedScore(adjustedHomeRating, awayRating)
	expectedAway := 1 - expectedHome

	// Draw factor approximation
	drawFactor := 0.26
	homeWin := expectedHome * (1 - drawFactor)
	awayWin := expectedAway * (1 - drawFactor)
	draw := drawFactor

	total := homeWin + draw + awayWin

	return ELOMatchProbabilities{
		HomeWin: (homeWin / total) * 100,
		Draw:    (draw / total) * 100,
		AwayWin: (awayWin / total) * 100,
	}
}

// ELOUpdateResult contains updated ratings for both teams.
type ELOUpdateResult struct {
	HomeRating ELORating
	AwayRating ELORating
}

// UpdateRatings calculates new ratings for both teams after a match.
func UpdateRatings(homeRating, awayRating float64, homeScore, awayScore int, kFactor, homeAdvantage float64) ELOUpdateResult {
	if kFactor == 0 {
		kFactor = KFactor
	}
	if homeAdvantage == 0 {
		homeAdvantage = HomeAdvantage
	}

	adjustedHomeRating := homeRating + homeAdvantage
	expectedHome := CalculateExpectedScore(adjustedHomeRating, awayRating)
	expectedAway := 1 - expectedHome

	var actualHome, actualAway float64
	if homeScore > awayScore {
		actualHome = 1
		actualAway = 0
	} else if homeScore < awayScore {
		actualHome = 0
		actualAway = 1
	} else {
		actualHome = 0.5
		actualAway = 0.5
	}

	// Goal difference multiplier
	goalDiff := abs(homeScore - awayScore)
	marginMultiplier := 1.0
	if goalDiff == 2 {
		marginMultiplier = 1.5
	} else if goalDiff > 2 {
		marginMultiplier = (11 + float64(goalDiff)) / 8
	}

	adjustedK := kFactor * marginMultiplier

	newHomeRating := CalculateNewRating(homeRating, expectedHome, actualHome, adjustedK)
	newAwayRating := CalculateNewRating(awayRating, expectedAway, actualAway, adjustedK)

	return ELOUpdateResult{
		HomeRating: ELORating{
			Rating: int(math.Round(newHomeRating)),
			Change: int(math.Round(newHomeRating - homeRating)),
		},
		AwayRating: ELORating{
			Rating: int(math.Round(newAwayRating)),
			Change: int(math.Round(newAwayRating - awayRating)),
		},
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// SimulateELOSeason simulates a season of matches and returns final ratings.
func SimulateELOSeason(teams map[string]float64, matches []MatchResult, kFactor float64) map[string]float64 {
	if kFactor == 0 {
		kFactor = KFactor
	}

	ratings := make(map[string]float64)
	for team, rating := range teams {
		ratings[team] = rating
	}

	for _, match := range matches {
		homeRating := ratings[match.HomeTeam]
		if homeRating == 0 {
			homeRating = BaseELO
		}
		awayRating := ratings[match.AwayTeam]
		if awayRating == 0 {
			awayRating = BaseELO
		}

		result := UpdateRatings(homeRating, awayRating, match.HomeScore, match.AwayScore, kFactor, HomeAdvantage)
		ratings[match.HomeTeam] = float64(result.HomeRating.Rating)
		ratings[match.AwayTeam] = float64(result.AwayRating.Rating)
	}

	return ratings
}

// GetInitialRating returns the default starting ELO rating.
func GetInitialRating() int {
	return BaseELO
}

// RatingToTier converts an ELO rating to a tier description.
func RatingToTier(rating int) string {
	switch {
	case rating >= 2000:
		return "Elite"
	case rating >= 1800:
		return "Strong"
	case rating >= 1600:
		return "Above Average"
	case rating >= 1400:
		return "Average"
	case rating >= 1200:
		return "Below Average"
	default:
		return "Weak"
	}
}

// PredictMatchOutcome combines ELO with other factors for prediction.
func PredictMatchOutcome(homeELO, awayELO float64, homeForm, awayForm float64) ELOMatchProbabilities {
	// Adjust ratings based on recent form (0-1 scale where 1 is perfect)
	formWeight := 0.1
	adjustedHome := homeELO * (1 + formWeight*(homeForm-0.5))
	adjustedAway := awayELO * (1 + formWeight*(awayForm-0.5))

	return CalculateELOMatchProbabilities(adjustedHome, adjustedAway, HomeAdvantage)
}
