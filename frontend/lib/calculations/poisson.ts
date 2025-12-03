import { PoissonPrediction } from '@/types/betting';

function factorial(n: number): number {
  if (n <= 1) return 1;
  let result = 1;
  for (let i = 2; i <= n; i++) {
    result *= i;
  }
  return result;
}

function poissonProbability(lambda: number, k: number): number {
  return (Math.pow(lambda, k) * Math.exp(-lambda)) / factorial(k);
}

export function calculatePoissonPrediction(
  homeGoalsAvg: number,
  homeConcededAvg: number,
  awayGoalsAvg: number,
  awayConcededAvg: number,
  leagueAvgGoals: number = 2.75
): PoissonPrediction {
  const homeAttackStrength = homeGoalsAvg / (leagueAvgGoals / 2);
  const homeDefenseStrength = homeConcededAvg / (leagueAvgGoals / 2);
  const awayAttackStrength = awayGoalsAvg / (leagueAvgGoals / 2);
  const awayDefenseStrength = awayConcededAvg / (leagueAvgGoals / 2);

  const expectedHomeGoals = homeAttackStrength * awayDefenseStrength * (leagueAvgGoals / 2);
  const expectedAwayGoals = awayAttackStrength * homeDefenseStrength * (leagueAvgGoals / 2);

  const maxGoals = 10;
  const scoreMatrix: number[][] = [];
  
  for (let i = 0; i <= maxGoals; i++) {
    scoreMatrix[i] = [];
    for (let j = 0; j <= maxGoals; j++) {
      scoreMatrix[i][j] = poissonProbability(expectedHomeGoals, i) * poissonProbability(expectedAwayGoals, j);
    }
  }

  let homeWinProb = 0;
  let drawProb = 0;
  let awayWinProb = 0;
  let over25Prob = 0;
  let bttsProb = 0;

  for (let i = 0; i <= maxGoals; i++) {
    for (let j = 0; j <= maxGoals; j++) {
      const prob = scoreMatrix[i][j];
      if (i > j) homeWinProb += prob;
      else if (i === j) drawProb += prob;
      else awayWinProb += prob;
      
      if (i + j > 2) over25Prob += prob;
      if (i > 0 && j > 0) bttsProb += prob;
    }
  }

  const mostLikelyScores: { score: string; probability: number }[] = [];
  for (let i = 0; i <= 5; i++) {
    for (let j = 0; j <= 5; j++) {
      mostLikelyScores.push({
        score: `${i}-${j}`,
        probability: scoreMatrix[i][j] * 100,
      });
    }
  }
  mostLikelyScores.sort((a, b) => b.probability - a.probability);

  return {
    homeGoals: expectedHomeGoals,
    awayGoals: expectedAwayGoals,
    homeWinProb: homeWinProb * 100,
    drawProb: drawProb * 100,
    awayWinProb: awayWinProb * 100,
    over25Prob: over25Prob * 100,
    under25Prob: (1 - over25Prob) * 100,
    bttsProb: bttsProb * 100,
    scoreMatrix,
    mostLikelyScores: mostLikelyScores.slice(0, 10),
  };
}

export function calculateMatchProbabilities(
  homeGoalsAvg: number,
  awayGoalsAvg: number
): { homeWin: number; draw: number; awayWin: number } {
  const prediction = calculatePoissonPrediction(
    homeGoalsAvg,
    2.5,
    awayGoalsAvg,
    2.5
  );
  
  return {
    homeWin: prediction.homeWinProb,
    draw: prediction.drawProb,
    awayWin: prediction.awayWinProb,
  };
}
