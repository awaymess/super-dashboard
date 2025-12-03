import { baseApi } from './baseApi';
import { Match, Bet, ValueBet, BettingStats } from '@/types/betting';

export const bettingApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getMatches: builder.query<Match[], { league?: string; date?: string }>({
      query: (params) => ({
        url: '/betting/matches',
        params,
      }),
      providesTags: ['Match'],
    }),
    getMatch: builder.query<Match, string>({
      query: (id) => `/betting/matches/${id}`,
      providesTags: (_result, _error, id) => [{ type: 'Match', id }],
    }),
    getBets: builder.query<Bet[], { status?: string; page?: number; limit?: number }>({
      query: (params) => ({
        url: '/betting/bets',
        params,
      }),
      providesTags: ['Bet'],
    }),
    placeBet: builder.mutation<Bet, Partial<Bet>>({
      query: (body) => ({
        url: '/betting/bets',
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Bet'],
    }),
    getValueBets: builder.query<ValueBet[], void>({
      query: () => '/betting/value-bets',
    }),
    getBettingStats: builder.query<BettingStats, void>({
      query: () => '/betting/stats',
    }),
  }),
});

export const {
  useGetMatchesQuery,
  useGetMatchQuery,
  useGetBetsQuery,
  usePlaceBetMutation,
  useGetValueBetsQuery,
  useGetBettingStatsQuery,
} = bettingApi;
