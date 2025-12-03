import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Portfolio, Position, Transaction, JournalEntry, TradeOrder, BacktestResult } from '@/types/paper-trading';

interface PaperTradingState {
  portfolio: Portfolio | null;
  positions: Position[];
  transactions: Transaction[];
  journalEntries: JournalEntry[];
  backtestResults: BacktestResult[];
  pendingOrder: TradeOrder | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: PaperTradingState = {
  portfolio: null,
  positions: [],
  transactions: [],
  journalEntries: [],
  backtestResults: [],
  pendingOrder: null,
  isLoading: false,
  error: null,
};

const paperTradingSlice = createSlice({
  name: 'paperTrading',
  initialState,
  reducers: {
    setPortfolio: (state, action: PayloadAction<Portfolio>) => {
      state.portfolio = action.payload;
      state.positions = action.payload.positions;
      state.transactions = action.payload.transactions;
    },
    updatePortfolioValue: (state, action: PayloadAction<{ totalValue: number; cashBalance: number }>) => {
      if (state.portfolio) {
        state.portfolio.totalValue = action.payload.totalValue;
        state.portfolio.cashBalance = action.payload.cashBalance;
      }
    },
    setPositions: (state, action: PayloadAction<Position[]>) => {
      state.positions = action.payload;
    },
    updatePosition: (state, action: PayloadAction<Position>) => {
      const index = state.positions.findIndex((p) => p.id === action.payload.id);
      if (index !== -1) {
        state.positions[index] = action.payload;
      }
    },
    addPosition: (state, action: PayloadAction<Position>) => {
      state.positions.push(action.payload);
    },
    removePosition: (state, action: PayloadAction<string>) => {
      state.positions = state.positions.filter((p) => p.id !== action.payload);
    },
    setTransactions: (state, action: PayloadAction<Transaction[]>) => {
      state.transactions = action.payload;
    },
    addTransaction: (state, action: PayloadAction<Transaction>) => {
      state.transactions.unshift(action.payload);
    },
    setJournalEntries: (state, action: PayloadAction<JournalEntry[]>) => {
      state.journalEntries = action.payload;
    },
    addJournalEntry: (state, action: PayloadAction<JournalEntry>) => {
      state.journalEntries.unshift(action.payload);
    },
    updateJournalEntry: (state, action: PayloadAction<JournalEntry>) => {
      const index = state.journalEntries.findIndex((e) => e.id === action.payload.id);
      if (index !== -1) {
        state.journalEntries[index] = action.payload;
      }
    },
    deleteJournalEntry: (state, action: PayloadAction<string>) => {
      state.journalEntries = state.journalEntries.filter((e) => e.id !== action.payload);
    },
    setBacktestResults: (state, action: PayloadAction<BacktestResult[]>) => {
      state.backtestResults = action.payload;
    },
    addBacktestResult: (state, action: PayloadAction<BacktestResult>) => {
      state.backtestResults.unshift(action.payload);
    },
    setPendingOrder: (state, action: PayloadAction<TradeOrder | null>) => {
      state.pendingOrder = action.payload;
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
  setPortfolio,
  updatePortfolioValue,
  setPositions,
  updatePosition,
  addPosition,
  removePosition,
  setTransactions,
  addTransaction,
  setJournalEntries,
  addJournalEntry,
  updateJournalEntry,
  deleteJournalEntry,
  setBacktestResults,
  addBacktestResult,
  setPendingOrder,
  setLoading,
  setError,
} = paperTradingSlice.actions;

export default paperTradingSlice.reducer;
