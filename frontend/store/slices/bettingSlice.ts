import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Match, Bet, BetSlipItem, BettingStats, ValueBet } from '@/types/betting';

interface BettingState {
  matches: Match[];
  selectedMatch: Match | null;
  betSlip: BetSlipItem[];
  bettingHistory: Bet[];
  valueBets: ValueBet[];
  stats: BettingStats;
  filters: {
    league: string;
    date: string;
    status: string;
  };
  isLoading: boolean;
  error: string | null;
}

const initialState: BettingState = {
  matches: [],
  selectedMatch: null,
  betSlip: [],
  bettingHistory: [],
  valueBets: [],
  stats: {
    totalBets: 0,
    wonBets: 0,
    lostBets: 0,
    voidBets: 0,
    totalStaked: 0,
    totalReturns: 0,
    profit: 0,
    roi: 0,
    avgOdds: 0,
    avgStake: 0,
    winRate: 0,
    currentStreak: 0,
    bestStreak: 0,
    worstStreak: 0,
  },
  filters: {
    league: '',
    date: '',
    status: '',
  },
  isLoading: false,
  error: null,
};

const bettingSlice = createSlice({
  name: 'betting',
  initialState,
  reducers: {
    setMatches: (state, action: PayloadAction<Match[]>) => {
      state.matches = action.payload;
    },
    selectMatch: (state, action: PayloadAction<Match | null>) => {
      state.selectedMatch = action.payload;
    },
    addToBetSlip: (state, action: PayloadAction<BetSlipItem>) => {
      const exists = state.betSlip.find(
        (item) => item.matchId === action.payload.matchId && item.betType === action.payload.betType
      );
      if (!exists) {
        state.betSlip.push(action.payload);
      }
    },
    removeFromBetSlip: (state, action: PayloadAction<{ matchId: string; betType: string }>) => {
      state.betSlip = state.betSlip.filter(
        (item) => !(item.matchId === action.payload.matchId && item.betType === action.payload.betType)
      );
    },
    updateBetSlipStake: (state, action: PayloadAction<{ matchId: string; betType: string; stake: number }>) => {
      const item = state.betSlip.find(
        (item) => item.matchId === action.payload.matchId && item.betType === action.payload.betType
      );
      if (item) {
        item.stake = action.payload.stake;
      }
    },
    clearBetSlip: (state) => {
      state.betSlip = [];
    },
    setBettingHistory: (state, action: PayloadAction<Bet[]>) => {
      state.bettingHistory = action.payload;
    },
    addBet: (state, action: PayloadAction<Bet>) => {
      state.bettingHistory.unshift(action.payload);
    },
    updateBet: (state, action: PayloadAction<Bet>) => {
      const index = state.bettingHistory.findIndex((bet) => bet.id === action.payload.id);
      if (index !== -1) {
        state.bettingHistory[index] = action.payload;
      }
    },
    setValueBets: (state, action: PayloadAction<ValueBet[]>) => {
      state.valueBets = action.payload;
    },
    setStats: (state, action: PayloadAction<BettingStats>) => {
      state.stats = action.payload;
    },
    setFilters: (state, action: PayloadAction<Partial<BettingState['filters']>>) => {
      state.filters = { ...state.filters, ...action.payload };
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
  },
});

export const {
  setMatches,
  selectMatch,
  addToBetSlip,
  removeFromBetSlip,
  updateBetSlipStake,
  clearBetSlip,
  setBettingHistory,
  addBet,
  updateBet,
  setValueBets,
  setStats,
  setFilters,
  setLoading,
  setError,
} = bettingSlice.actions;

export default bettingSlice.reducer;
