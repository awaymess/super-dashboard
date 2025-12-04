# Calculations Package

This package provides mathematical models for betting analytics and stock valuation in Go. These functions mirror the frontend TypeScript implementations in `frontend/lib/calculations`.

## Overview

The package includes the following calculation modules:

- **Poisson Distribution** - Football/soccer match outcome predictions
- **Kelly Criterion** - Optimal bet sizing with fractional variants
- **ELO Rating** - Team strength rating system
- **DCF Valuation** - Discounted Cash Flow stock valuation
- **Graham Formula** - Benjamin Graham intrinsic value calculations
- **Monte Carlo** - Simulation for bankroll growth and risk analysis

## Installation

```go
import "github.com/superdashboard/backend/lib/calculations"
```

## Usage Examples

### Poisson Distribution

Calculate match outcome probabilities using the Poisson distribution:

```go
// Calculate probability of a team scoring exactly 2 goals (lambda = 1.5)
prob := calculations.PoissonProbability(1.5, 2)
// Result: ~0.251

// Full match prediction
// Parameters in order:
//   homeGoalsAvg      - Average goals scored by home team
//   homeConcededAvg   - Average goals conceded by home team
//   awayGoalsAvg      - Average goals scored by away team
//   awayConcededAvg   - Average goals conceded by away team
//   leagueAvgGoals    - League average total goals per match
prediction := calculations.CalculatePoissonPrediction(
    1.8,   // homeGoalsAvg: Home team scores 1.8 goals/game on average
    1.0,   // homeConcededAvg: Home team concedes 1.0 goals/game
    1.5,   // awayGoalsAvg: Away team scores 1.5 goals/game
    1.2,   // awayConcededAvg: Away team concedes 1.2 goals/game
    2.75,  // leagueAvgGoals: League averages 2.75 total goals/game
)
// prediction.HomeWinProb, prediction.DrawProb, prediction.AwayWinProb
// prediction.Over25Prob, prediction.BTTSProb
// prediction.MostLikelyScores (top 10 score probabilities)
```

### Kelly Criterion

Calculate optimal bet sizing:

```go
// Full Kelly
kelly := calculations.CalculateKelly(
    60.0,   // Win probability (%)
    2.0,    // Decimal odds
    1000.0, // Bankroll
    1.0,    // Fraction (1.0 = full Kelly)
)
// kelly.Stake         = optimal stake amount
// kelly.HalfKelly     = half Kelly stake
// kelly.QuarterKelly  = quarter Kelly stake
// kelly.ExpectedValue = expected value (%)

// Using individual functions
fullKelly := calculations.CalculateFullKelly(0.6, 2.0)    // Returns 0.2 (20% of bankroll)
halfKelly := calculations.CalculateHalfKelly(0.6, 2.0)    // Returns 0.1 (10%)
quarterKelly := calculations.CalculateQuarterKelly(0.6, 2.0) // Returns 0.05 (5%)

// Value bet detection
valueBet := calculations.DetectValueBet(
    60.0, // True probability (%)
    2.0,  // Bookmaker odds
    5.0,  // Value threshold (%)
)
// valueBet.IsValueBet, valueBet.Value, valueBet.Recommendation
```

### ELO Rating

Calculate and update ELO ratings:

```go
// Calculate expected score
expected := calculations.CalculateExpectedScore(1600, 1500)
// Result: ~0.64 (64% expected win probability for higher rated player)

// Update ratings after a match
result := calculations.UpdateRatings(
    1500,  // Home rating
    1500,  // Away rating
    2,     // Home score
    1,     // Away score
    32,    // K-factor (0 for default)
    100,   // Home advantage (0 for default)
)
// result.HomeRating.Rating, result.HomeRating.Change
// result.AwayRating.Rating, result.AwayRating.Change

// Match probability prediction
probs := calculations.CalculateELOMatchProbabilities(1600, 1400, 100)
// probs.HomeWin, probs.Draw, probs.AwayWin (as percentages)

// Simulate a season
teams := map[string]float64{"TeamA": 1500, "TeamB": 1500, "TeamC": 1500}
matches := []calculations.MatchResult{
    {HomeTeam: "TeamA", AwayTeam: "TeamB", HomeScore: 2, AwayScore: 1},
    {HomeTeam: "TeamB", AwayTeam: "TeamC", HomeScore: 1, AwayScore: 1},
}
finalRatings := calculations.SimulateELOSeason(teams, matches, 32)
```

### DCF Valuation

Perform Discounted Cash Flow analysis:

```go
inputs := calculations.DCFInputs{
    FreeCashFlow:       100000000, // $100M
    GrowthRate:         10,        // 10% annual growth
    TerminalGrowthRate: 3,         // 3% terminal growth
    DiscountRate:       10,        // 10% WACC
    Years:              5,         // 5-year projection
    SharesOutstanding:  10000000,  // 10M shares
}

result := calculations.CalculateDCF(inputs)
// result.IntrinsicValue        = total enterprise value
// result.PerShareValue         = fair value per share
// result.PresentValueOfCashFlows
// result.TerminalValue
// result.ProjectedCashFlows    = year-by-year projections

// Calculate WACC
wacc := calculations.CalculateWACC(60, 40, 12, 6, 25)
// 60% equity at 12%, 40% debt at 6%, 25% tax rate

// Calculate margin of safety
mos := calculations.CalculateMarginOfSafety(100, 70)
// Result: 30% (stock worth $100, trading at $70)
```

### Graham Formula

Benjamin Graham valuation methods:

```go
// Graham Number: sqrt(22.5 × EPS × Book Value)
grahamNumber := calculations.CalculateGrahamNumber(5.0, 40.0)
// Result: ~67.08

// Modified Graham Value: EPS × (8.5 + 2g) × 4.4 / Y
modifiedValue := calculations.CalculateModifiedGrahamValue(5.0, 10.0, 4.4)
// Result: 142.5

// Full analysis
inputs := calculations.GrahamInputs{
    EPS:               5.0,
    BookValuePerShare: 40.0,
    CurrentPrice:      50.0,
    GrowthRate:        10.0,
    AAAYield:          4.4,
}
result := calculations.CalculateGrahamAnalysis(inputs)
// result.GrahamNumber, result.ModifiedGrahamValue
// result.MarginOfSafety, result.Rating
// result.Analysis.PERatio, result.Analysis.PBRatio
// result.Analysis.IsDefensive, result.Analysis.IsEnterprising

// Screen for defensive stocks
stocks := []calculations.StockForScreening{
    {Symbol: "GOOD", EPS: 5, BookValue: 40, Price: 50},
    {Symbol: "HIGHPE", EPS: 5, BookValue: 40, Price: 100},
}
defensive := calculations.ScreenDefensiveStocks(stocks)
// Returns: ["GOOD"]

// Net Current Asset Value (NCAV)
// NCAV = Current Assets - Total Liabilities
// Parameters: currentAssets, totalLiabilities, sharesOutstanding
ncav := calculations.CalculateNCAV(
    100000000, // currentAssets: $100M in current assets
    60000000,  // totalLiabilities: $60M in total liabilities
    10000000,  // sharesOutstanding: 10M shares
)
// ncav.NCAV = $40M, ncav.NCAVPerShare = $4

// Check for net-net bargains (market cap < 2/3 of NCAV)
// Parameters: currentAssets, totalLiabilities, marketCap
isNetNet := calculations.IsNetNet(
    100000000, // currentAssets: $100M
    60000000,  // totalLiabilities: $60M
    20000000,  // marketCap: $20M
)
// Returns true because $20M < (2/3 * $40M NCAV)
```

### Monte Carlo Simulation

Run simulations for investment and betting scenarios:

```go
// Investment portfolio simulation
config := calculations.MonteCarloConfig{
    Simulations:      10000,
    InitialValue:     100000,
    ExpectedReturn:   8,    // 8% annual
    Volatility:       15,   // 15% std dev
    TimeHorizonYears: 1,
    Seed:             42,   // 0 for random
}

result := calculations.RunMonteCarloSimulation(config)
// result.Mean, result.Median, result.StandardDeviation
// result.Percentile5, result.Percentile25, result.Percentile75, result.Percentile95
// result.MinValue, result.MaxValue
// result.ProbabilityOfLoss, result.ProbabilityOfGain

// Betting simulation
bettingConfig := calculations.BettingMonteCarloConfig{
    Simulations:     10000,
    InitialBankroll: 1000,
    NumBets:         100,
    WinProbability:  55,   // 55% win rate
    AverageOdds:     2.0,
    StakePercent:    2,    // 2% per bet
    Seed:            42,
}

bettingResult := calculations.RunBettingMonteCarlo(bettingConfig)
// bettingResult.AvgFinalBankroll
// bettingResult.AvgMaxDrawdown
// bettingResult.RuinProbability
// bettingResult.DoubleProbability

// Calculate Sharpe ratio
returns := []float64{10, 15, 5, 20, 8}
sharpe := calculations.SharpeRatio(returns, 2.0) // 2% risk-free rate

// Calculate max drawdown
values := []float64{100, 110, 90, 95, 105}
maxDD := calculations.MaxDrawdown(values)
// Result: ~18.18%

// Value at Risk
var95 := calculations.ValueAtRisk(returns, 95)
```

## Running Tests

```bash
cd backend
go test -v ./lib/calculations/...
```

All functions include comprehensive unit tests with deterministic examples for reproducibility.

## Constants

```go
const (
    BaseELO       = 1500  // Starting ELO rating
    KFactor       = 32    // Default K-factor for ELO updates
    HomeAdvantage = 100   // Default home advantage in ELO points
)
```

## Error Handling

Most functions handle edge cases gracefully:
- Negative inputs return 0 or empty results
- Division by zero is protected
- Invalid probabilities are clamped to valid ranges

## Thread Safety

All functions are stateless and thread-safe. The Monte Carlo simulations use their own random sources.
