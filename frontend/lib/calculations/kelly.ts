import { KellyResult } from '@/types/betting';

export function calculateKelly(
  probability: number,
  odds: number,
  bankroll: number = 1000,
  fractionalKelly: number = 1
): KellyResult {
  const p = probability / 100;
  const q = 1 - p;
  const b = odds - 1;

  let kellyFraction = (b * p - q) / b;
  
  if (kellyFraction < 0) {
    kellyFraction = 0;
  }

  const fullKellyStake = kellyFraction * bankroll * fractionalKelly;
  const edge = (p * odds) - 1;
  const expectedValue = edge * 100;

  return {
    stake: Math.max(0, fullKellyStake),
    halfKelly: Math.max(0, fullKellyStake * 0.5),
    quarterKelly: Math.max(0, fullKellyStake * 0.25),
    expectedValue,
    edge: edge * 100,
  };
}

export function calculateOptimalStake(
  probability: number,
  odds: number,
  bankroll: number,
  maxStakePercent: number = 5
): number {
  const kelly = calculateKelly(probability, odds, bankroll);
  const maxStake = bankroll * (maxStakePercent / 100);
  return Math.min(kelly.halfKelly, maxStake);
}

export function calculateValue(
  fairProbability: number,
  bookmakerOdds: number
): { value: number; isValueBet: boolean; expectedValue: number } {
  const fairOdds = 100 / fairProbability;
  const value = ((bookmakerOdds / fairOdds) - 1) * 100;
  const expectedValue = (fairProbability / 100) * (bookmakerOdds - 1) - (1 - fairProbability / 100);
  
  return {
    value,
    isValueBet: value > 0,
    expectedValue: expectedValue * 100,
  };
}

export function calculateKellyWithEdge(
  edge: number,
  odds: number,
  bankroll: number
): number {
  if (edge <= 0) return 0;
  const b = odds - 1;
  const p = (edge / 100) + (1 / odds);
  const q = 1 - p;
  const kelly = (b * p - q) / b;
  return Math.max(0, kelly * bankroll);
}

export function simulateKellyGrowth(
  initialBankroll: number,
  bets: { probability: number; odds: number; won: boolean }[],
  fraction: number = 0.5
): { finalBankroll: number; growth: number; maxDrawdown: number } {
  let bankroll = initialBankroll;
  let maxBankroll = initialBankroll;
  let maxDrawdown = 0;

  for (const bet of bets) {
    const stake = calculateKelly(bet.probability, bet.odds, bankroll, fraction).stake;
    
    if (bet.won) {
      bankroll += stake * (bet.odds - 1);
    } else {
      bankroll -= stake;
    }

    if (bankroll > maxBankroll) {
      maxBankroll = bankroll;
    }

    const drawdown = ((maxBankroll - bankroll) / maxBankroll) * 100;
    if (drawdown > maxDrawdown) {
      maxDrawdown = drawdown;
    }
  }

  return {
    finalBankroll: bankroll,
    growth: ((bankroll - initialBankroll) / initialBankroll) * 100,
    maxDrawdown,
  };
}

/**
 * Calculate full Kelly criterion stake as a fraction
 * @param probability - Win probability (0 to 1)
 * @param odds - Decimal odds
 * @returns Kelly stake as a fraction of bankroll
 */
export function calculateFullKelly(probability: number, odds: number): number {
  if (probability <= 0 || probability >= 1 || odds <= 1) return 0;
  const p = probability;
  const q = 1 - p;
  const b = odds - 1;
  const kelly = (b * p - q) / b;
  return Math.max(0, kelly);
}

/**
 * Calculate half Kelly criterion stake as a fraction
 * @param probability - Win probability (0 to 1)
 * @param odds - Decimal odds
 * @returns Half Kelly stake as a fraction of bankroll
 */
export function calculateHalfKelly(probability: number, odds: number): number {
  return calculateFullKelly(probability, odds) * 0.5;
}

/**
 * Calculate quarter Kelly criterion stake as a fraction
 * @param probability - Win probability (0 to 1)
 * @param odds - Decimal odds
 * @returns Quarter Kelly stake as a fraction of bankroll
 */
export function calculateQuarterKelly(probability: number, odds: number): number {
  return calculateFullKelly(probability, odds) * 0.25;
}

/**
 * Calculate implied probability from decimal odds
 * @param decimalOdds - Decimal odds (e.g., 2.00)
 * @returns Implied probability as percentage (0-100)
 */
export function calculateImpliedProbability(decimalOdds: number): number {
  if (decimalOdds <= 1) return 0;
  return (1 / decimalOdds) * 100;
}

/**
 * Convert probability to fair odds
 * @param probability - Probability as percentage (0-100)
 * @returns Fair decimal odds
 */
export function probabilityToOdds(probability: number): number {
  if (probability <= 0 || probability >= 100) return 0;
  return 100 / probability;
}

/**
 * Calculate value bet with threshold detection
 * @param trueProbability - True probability as percentage (0-100)
 * @param bookmakerOdds - Bookmaker decimal odds
 * @param valueThreshold - Minimum value percentage to consider (default 5%)
 * @returns Value bet analysis result
 */
export function detectValueBet(
  trueProbability: number,
  bookmakerOdds: number,
  valueThreshold: number = 5
): {
  impliedProbability: number;
  value: number;
  isValueBet: boolean;
  isHighValue: boolean;
  expectedValue: number;
  recommendation: 'skip' | 'bet' | 'strong_bet';
} {
  const impliedProbability = calculateImpliedProbability(bookmakerOdds);
  const value = trueProbability - impliedProbability;
  const expectedValue = (trueProbability / 100) * (bookmakerOdds - 1) - (1 - trueProbability / 100);
  
  let recommendation: 'skip' | 'bet' | 'strong_bet' = 'skip';
  if (value > 10) {
    recommendation = 'strong_bet';
  } else if (value > valueThreshold) {
    recommendation = 'bet';
  }

  return {
    impliedProbability,
    value,
    isValueBet: value > valueThreshold,
    isHighValue: value > 10,
    expectedValue: expectedValue * 100,
    recommendation,
  };
}

/**
 * Bayesian probability update based on new evidence
 * @param priorProbability - Prior probability (0-1)
 * @param likelihood - Probability of evidence given hypothesis (0-1)
 * @param evidenceProbability - Overall probability of evidence (0-1)
 * @returns Updated (posterior) probability
 */
export function bayesianUpdate(
  priorProbability: number,
  likelihood: number,
  evidenceProbability: number
): number {
  if (evidenceProbability === 0) return priorProbability;
  return (likelihood * priorProbability) / evidenceProbability;
}

/**
 * Calculate weighted average probability from multiple models
 * @param probabilities - Array of probability estimates (0-100)
 * @param weights - Array of weights for each model (should sum to 1)
 * @returns Weighted average probability
 */
export function calculateWeightedProbability(
  probabilities: number[],
  weights: number[]
): number {
  if (probabilities.length !== weights.length) {
    throw new Error(`Probabilities array length (${probabilities.length}) must match weights array length (${weights.length})`);
  }
  
  const totalWeight = weights.reduce((sum, w) => sum + w, 0);
  if (Math.abs(totalWeight - 1) > 0.001) {
    // Normalize weights if they don't sum to 1
    const normalizedWeights = weights.map(w => w / totalWeight);
    return probabilities.reduce((sum, prob, i) => sum + prob * normalizedWeights[i], 0);
  }
  
  return probabilities.reduce((sum, prob, i) => sum + prob * weights[i], 0);
}

/**
 * Calculate ensemble probability from multiple models
 * Uses equal weights by default
 * @param poissonProb - Probability from Poisson model
 * @param eloProb - Probability from ELO model
 * @param statProb - Probability from statistics-based model
 * @param xgProb - Probability from xG model (optional)
 * @returns Ensemble probability
 */
export function calculateEnsembleProbability(
  poissonProb: number,
  eloProb: number,
  statProb: number,
  xgProb?: number
): number {
  const probs = [poissonProb, eloProb, statProb];
  if (xgProb !== undefined) {
    probs.push(xgProb);
  }
  const equalWeight = 1 / probs.length;
  const weights = probs.map(() => equalWeight);
  return calculateWeightedProbability(probs, weights);
}

/**
 * Find arbitrage opportunities between bookmakers
 * @param odds - Array of odds from different bookmakers [home, draw, away][]
 * @returns Arbitrage analysis
 */
export function findArbitrage(
  homeOdds: number[],
  drawOdds: number[],
  awayOdds: number[]
): {
  isArbitrage: boolean;
  margin: number;
  bestHome: { index: number; odds: number };
  bestDraw: { index: number; odds: number };
  bestAway: { index: number; odds: number };
  stakes?: { home: number; draw: number; away: number };
} {
  const bestHome = { index: 0, odds: Math.max(...homeOdds) };
  bestHome.index = homeOdds.indexOf(bestHome.odds);
  
  const bestDraw = { index: 0, odds: Math.max(...drawOdds) };
  bestDraw.index = drawOdds.indexOf(bestDraw.odds);
  
  const bestAway = { index: 0, odds: Math.max(...awayOdds) };
  bestAway.index = awayOdds.indexOf(bestAway.odds);
  
  const totalImplied = (1 / bestHome.odds) + (1 / bestDraw.odds) + (1 / bestAway.odds);
  const margin = (1 - totalImplied) * 100;
  const isArbitrage = totalImplied < 1;
  
  let stakes;
  if (isArbitrage) {
    const total = 100;
    stakes = {
      home: (total / bestHome.odds) / totalImplied,
      draw: (total / bestDraw.odds) / totalImplied,
      away: (total / bestAway.odds) / totalImplied,
    };
  }
  
  return {
    isArbitrage,
    margin,
    bestHome,
    bestDraw,
    bestAway,
    stakes,
  };
}
