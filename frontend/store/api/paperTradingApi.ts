import { baseApi } from './baseApi';
import { Portfolio, Position, Transaction, TradeOrder, BacktestConfig, BacktestResult, JournalEntry, LeaderboardEntry } from '@/types/paper-trading';

export const paperTradingApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getPortfolio: builder.query<Portfolio, void>({
      query: () => '/paper-trading/portfolio',
      providesTags: ['Portfolio'],
    }),
    getPositions: builder.query<Position[], void>({
      query: () => '/paper-trading/positions',
      providesTags: ['Position'],
    }),
    getTransactions: builder.query<Transaction[], { page?: number; limit?: number }>({
      query: (params) => ({
        url: '/paper-trading/transactions',
        params,
      }),
      providesTags: ['Transaction'],
    }),
    executeTrade: builder.mutation<Transaction, TradeOrder>({
      query: (order) => ({
        url: '/paper-trading/trade',
        method: 'POST',
        body: order,
      }),
      invalidatesTags: ['Portfolio', 'Position', 'Transaction'],
    }),
    createLimitOrder: builder.mutation<{ orderId: string }, TradeOrder>({
      query: (order) => ({
        url: '/paper-trading/orders',
        method: 'POST',
        body: order,
      }),
    }),
    cancelOrder: builder.mutation<void, string>({
      query: (orderId) => ({
        url: `/paper-trading/orders/${orderId}`,
        method: 'DELETE',
      }),
    }),
    getJournalEntries: builder.query<JournalEntry[], void>({
      query: () => '/paper-trading/journal',
    }),
    createJournalEntry: builder.mutation<JournalEntry, Partial<JournalEntry>>({
      query: (entry) => ({
        url: '/paper-trading/journal',
        method: 'POST',
        body: entry,
      }),
    }),
    updateJournalEntry: builder.mutation<JournalEntry, { id: string; entry: Partial<JournalEntry> }>({
      query: ({ id, entry }) => ({
        url: `/paper-trading/journal/${id}`,
        method: 'PUT',
        body: entry,
      }),
    }),
    runBacktest: builder.mutation<BacktestResult, BacktestConfig>({
      query: (config) => ({
        url: '/paper-trading/backtest',
        method: 'POST',
        body: config,
      }),
    }),
    getBacktestResults: builder.query<BacktestResult[], void>({
      query: () => '/paper-trading/backtest/results',
    }),
    getLeaderboard: builder.query<LeaderboardEntry[], { period?: 'daily' | 'weekly' | 'monthly' | 'all' }>({
      query: (params) => ({
        url: '/paper-trading/leaderboard',
        params,
      }),
    }),
    resetPortfolio: builder.mutation<Portfolio, { initialBalance?: number }>({
      query: (body) => ({
        url: '/paper-trading/reset',
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Portfolio', 'Position', 'Transaction'],
    }),
  }),
});

export const {
  useGetPortfolioQuery,
  useGetPositionsQuery,
  useGetTransactionsQuery,
  useExecuteTradeMutation,
  useCreateLimitOrderMutation,
  useCancelOrderMutation,
  useGetJournalEntriesQuery,
  useCreateJournalEntryMutation,
  useUpdateJournalEntryMutation,
  useRunBacktestMutation,
  useGetBacktestResultsQuery,
  useGetLeaderboardQuery,
  useResetPortfolioMutation,
} = paperTradingApi;
