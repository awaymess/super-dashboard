const BASE_ELO = 1500;
const K_FACTOR = 32;
const HOME_ADVANTAGE = 100;

export interface EloRating {
  rating: number;
  change: number;
}

export interface MatchResult {
  homeTeam: string;
  awayTeam: string;
  homeScore: number;
  awayScore: number;
}

export function calculateExpectedScore(ratingA: number, ratingB: number): number {
  return 1 / (1 + Math.pow(10, (ratingB - ratingA) / 400));
}

export function calculateNewRating(
  rating: number,
  expectedScore: number,
  actualScore: number,
  kFactor: number = K_FACTOR
): number {
  return rating + kFactor * (actualScore - expectedScore);
}

export function calculateMatchProbabilities(
  homeRating: number,
  awayRating: number,
  homeAdvantage: number = HOME_ADVANTAGE
): { homeWin: number; draw: number; awayWin: number } {
  const adjustedHomeRating = homeRating + homeAdvantage;
  const expectedHome = calculateExpectedScore(adjustedHomeRating, awayRating);
  const expectedAway = 1 - expectedHome;

  const drawFactor = 0.26;
  const homeWin = expectedHome * (1 - drawFactor);
  const awayWin = expectedAway * (1 - drawFactor);
  const draw = drawFactor;

  const total = homeWin + draw + awayWin;
  
  return {
    homeWin: (homeWin / total) * 100,
    draw: (draw / total) * 100,
    awayWin: (awayWin / total) * 100,
  };
}

export function updateRatings(
  homeRating: number,
  awayRating: number,
  homeScore: number,
  awayScore: number,
  kFactor: number = K_FACTOR,
  homeAdvantage: number = HOME_ADVANTAGE
): { homeRating: EloRating; awayRating: EloRating } {
  const adjustedHomeRating = homeRating + homeAdvantage;
  const expectedHome = calculateExpectedScore(adjustedHomeRating, awayRating);
  const expectedAway = 1 - expectedHome;

  let actualHome: number;
  let actualAway: number;

  if (homeScore > awayScore) {
    actualHome = 1;
    actualAway = 0;
  } else if (homeScore < awayScore) {
    actualHome = 0;
    actualAway = 1;
  } else {
    actualHome = 0.5;
    actualAway = 0.5;
  }

  const goalDiff = Math.abs(homeScore - awayScore);
  const marginMultiplier = goalDiff <= 1 ? 1 : goalDiff === 2 ? 1.5 : (11 + goalDiff) / 8;

  const adjustedK = kFactor * marginMultiplier;

  const newHomeRating = calculateNewRating(homeRating, expectedHome, actualHome, adjustedK);
  const newAwayRating = calculateNewRating(awayRating, expectedAway, actualAway, adjustedK);

  return {
    homeRating: {
      rating: Math.round(newHomeRating),
      change: Math.round(newHomeRating - homeRating),
    },
    awayRating: {
      rating: Math.round(newAwayRating),
      change: Math.round(newAwayRating - awayRating),
    },
  };
}

export function simulateEloSeason(
  teams: Record<string, number>,
  matches: MatchResult[],
  kFactor: number = K_FACTOR
): Record<string, number> {
  const ratings = { ...teams };

  for (const match of matches) {
    const result = updateRatings(
      ratings[match.homeTeam] || BASE_ELO,
      ratings[match.awayTeam] || BASE_ELO,
      match.homeScore,
      match.awayScore,
      kFactor
    );

    ratings[match.homeTeam] = result.homeRating.rating;
    ratings[match.awayTeam] = result.awayRating.rating;
  }

  return ratings;
}

export function getInitialRating(): number {
  return BASE_ELO;
}

export function ratingToTier(rating: number): string {
  if (rating >= 2000) return 'Elite';
  if (rating >= 1800) return 'Strong';
  if (rating >= 1600) return 'Above Average';
  if (rating >= 1400) return 'Average';
  if (rating >= 1200) return 'Below Average';
  return 'Weak';
}
