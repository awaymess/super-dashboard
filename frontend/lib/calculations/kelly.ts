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
