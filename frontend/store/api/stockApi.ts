import { baseApi } from './baseApi';
import { Stock, StockQuote, StockCandle, StockNews, AnalystRating, TechnicalIndicators } from '@/types/stocks';

export const stockApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getStocks: builder.query<Stock[], { sector?: string; exchange?: string }>({
      query: (params) => ({
        url: '/stocks',
        params,
      }),
      providesTags: ['Stock'],
    }),
    getStock: builder.query<Stock, string>({
      query: (symbol) => `/stocks/${symbol}`,
      providesTags: (_result, _error, symbol) => [{ type: 'Stock', id: symbol }],
    }),
    getQuote: builder.query<StockQuote, string>({
      query: (symbol) => `/stocks/${symbol}/quote`,
    }),
    getCandles: builder.query<StockCandle[], { symbol: string; interval: string; from?: string; to?: string }>({
      query: ({ symbol, ...params }) => ({
        url: `/stocks/${symbol}/candles`,
        params,
      }),
    }),
    getTechnicals: builder.query<TechnicalIndicators, string>({
      query: (symbol) => `/stocks/${symbol}/technicals`,
    }),
    getNews: builder.query<StockNews[], { symbols?: string[]; limit?: number }>({
      query: (params) => ({
        url: '/stocks/news',
        params: {
          symbols: params.symbols?.join(','),
          limit: params.limit,
        },
      }),
    }),
    getAnalystRatings: builder.query<AnalystRating, string>({
      query: (symbol) => `/stocks/${symbol}/analyst-ratings`,
    }),
    searchStocks: builder.query<Stock[], string>({
      query: (query) => ({
        url: '/stocks/search',
        params: { q: query },
      }),
    }),
  }),
});

export const {
  useGetStocksQuery,
  useGetStockQuery,
  useGetQuoteQuery,
  useGetCandlesQuery,
  useGetTechnicalsQuery,
  useGetNewsQuery,
  useGetAnalystRatingsQuery,
  useSearchStocksQuery,
} = stockApi;
