export interface Team {
  id: string;
  name: string;
  shortName: string;
  logo?: string;
  form: string[];
  homeForm?: string[];
  awayForm?: string[];
  goalsScoredAvg: number;
  goalsConcededAvg: number;
  cleanSheetPct: number;
}

export interface Match {
  id: string;
  homeTeam: Team;
  awayTeam: Team;
  league: string;
  leagueCountry: string;
  date: string;
  time: string;
  status: 'scheduled' | 'live' | 'finished' | 'postponed';
  homeScore?: number;
  awayScore?: number;
  odds: MatchOdds;
  stats?: MatchStats;
  h2h?: H2HRecord[];
}

export interface MatchOdds {
  home: number;
  draw: number;
  away: number;
  over25: number;
  under25: number;
  btts: number;
  bttsNo: number;
  homeDouble: number;
  awayDouble: number;
  correctScores?: Record<string, number>;
}

export interface MatchStats {
  possession: { home: number; away: number };
  shots: { home: number; away: number };
  shotsOnTarget: { home: number; away: number };
  corners: { home: number; away: number };
  fouls: { home: number; away: number };
  yellowCards: { home: number; away: number };
  redCards: { home: number; away: number };
}

export interface H2HRecord {
  date: string;
  homeTeam: string;
  awayTeam: string;
  homeScore: number;
  awayScore: number;
  competition: string;
}

export interface ValueBet {
  matchId: string;
  match: Match;
  betType: string;
  bookmakerOdds: number;
  fairOdds: number;
  value: number;
  confidence: number;
  kellyStake: number;
  expectedValue: number;
}

export interface Bet {
  id: string;
  matchId: string;
  match: Match;
  betType: string;
  odds: number;
  stake: number;
  potentialWin: number;
  status: 'pending' | 'won' | 'lost' | 'void';
  placedAt: string;
  settledAt?: string;
  profit?: number;
}

export interface BetSlip {
  bets: BetSlipItem[];
  totalStake: number;
  totalOdds: number;
  potentialWin: number;
}

export interface BetSlipItem {
  matchId: string;
  match: Match;
  betType: string;
  odds: number;
  stake: number;
}

export interface PoissonPrediction {
  homeGoals: number;
  awayGoals: number;
  homeWinProb: number;
  drawProb: number;
  awayWinProb: number;
  over25Prob: number;
  under25Prob: number;
  bttsProb: number;
  scoreMatrix: number[][];
  mostLikelyScores: { score: string; probability: number }[];
}

export interface KellyResult {
  stake: number;
  halfKelly: number;
  quarterKelly: number;
  expectedValue: number;
  edge: number;
}

export interface BettingStats {
  totalBets: number;
  wonBets: number;
  lostBets: number;
  voidBets: number;
  totalStaked: number;
  totalReturns: number;
  profit: number;
  roi: number;
  avgOdds: number;
  avgStake: number;
  winRate: number;
  currentStreak: number;
  bestStreak: number;
  worstStreak: number;
}
