import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Stock, StockQuote, Watchlist, StockNews, ScreenerCriteria } from '@/types/stocks';

interface StockState {
  stocks: Stock[];
  selectedStock: Stock | null;
  quote: StockQuote | null;
  watchlists: Watchlist[];
  activeWatchlist: Watchlist | null;
  news: StockNews[];
  screenerCriteria: ScreenerCriteria;
  screenerResults: Stock[];
  isLoading: boolean;
  error: string | null;
}

const initialState: StockState = {
  stocks: [],
  selectedStock: null,
  quote: null,
  watchlists: [
    {
      id: 'default',
      name: 'My Watchlist',
      symbols: ['AAPL', 'MSFT', 'GOOGL', 'AMZN', 'NVDA'],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ],
  activeWatchlist: null,
  news: [],
  screenerCriteria: {},
  screenerResults: [],
  isLoading: false,
  error: null,
};

const stockSlice = createSlice({
  name: 'stocks',
  initialState,
  reducers: {
    setStocks: (state, action: PayloadAction<Stock[]>) => {
      state.stocks = action.payload;
    },
    selectStock: (state, action: PayloadAction<Stock | null>) => {
      state.selectedStock = action.payload;
    },
    setQuote: (state, action: PayloadAction<StockQuote | null>) => {
      state.quote = action.payload;
    },
    setWatchlists: (state, action: PayloadAction<Watchlist[]>) => {
      state.watchlists = action.payload;
    },
    setActiveWatchlist: (state, action: PayloadAction<Watchlist | null>) => {
      state.activeWatchlist = action.payload;
    },
    createWatchlist: (state, action: PayloadAction<{ name: string; symbols?: string[] }>) => {
      const newWatchlist: Watchlist = {
        id: `watchlist-${Date.now()}`,
        name: action.payload.name,
        symbols: action.payload.symbols || [],
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
      state.watchlists.push(newWatchlist);
    },
    deleteWatchlist: (state, action: PayloadAction<string>) => {
      state.watchlists = state.watchlists.filter((w) => w.id !== action.payload);
      if (state.activeWatchlist?.id === action.payload) {
        state.activeWatchlist = null;
      }
    },
    addToWatchlist: (state, action: PayloadAction<{ watchlistId: string; symbol: string }>) => {
      const watchlist = state.watchlists.find((w) => w.id === action.payload.watchlistId);
      if (watchlist && !watchlist.symbols.includes(action.payload.symbol)) {
        watchlist.symbols.push(action.payload.symbol);
        watchlist.updatedAt = new Date().toISOString();
      }
    },
    removeFromWatchlist: (state, action: PayloadAction<{ watchlistId: string; symbol: string }>) => {
      const watchlist = state.watchlists.find((w) => w.id === action.payload.watchlistId);
      if (watchlist) {
        watchlist.symbols = watchlist.symbols.filter((s) => s !== action.payload.symbol);
        watchlist.updatedAt = new Date().toISOString();
      }
    },
    setNews: (state, action: PayloadAction<StockNews[]>) => {
      state.news = action.payload;
    },
    setScreenerCriteria: (state, action: PayloadAction<ScreenerCriteria>) => {
      state.screenerCriteria = action.payload;
    },
    setScreenerResults: (state, action: PayloadAction<Stock[]>) => {
      state.screenerResults = action.payload;
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
  setStocks,
  selectStock,
  setQuote,
  setWatchlists,
  setActiveWatchlist,
  createWatchlist,
  deleteWatchlist,
  addToWatchlist,
  removeFromWatchlist,
  setNews,
  setScreenerCriteria,
  setScreenerResults,
  setLoading,
  setError,
} = stockSlice.actions;

export default stockSlice.reducer;
