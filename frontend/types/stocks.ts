export interface Stock {
  symbol: string;
  name: string;
  exchange: string;
  sector: string;
  industry: string;
  marketCap: number;
  price: number;
  change: number;
  changePercent: number;
  volume: number;
  avgVolume: number;
  high52Week: number;
  low52Week: number;
  pe: number;
  eps: number;
  dividend: number;
  dividendYield: number;
  beta: number;
}

export interface StockQuote {
  symbol: string;
  price: number;
  change: number;
  changePercent: number;
  open: number;
  high: number;
  low: number;
  previousClose: number;
  volume: number;
  timestamp: string;
}

export interface StockCandle {
  date: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
}

export interface TechnicalIndicators {
  rsi: number;
  rsiSignal: 'oversold' | 'neutral' | 'overbought';
  macd: {
    value: number;
    signal: number;
    histogram: number;
    trend: 'bullish' | 'bearish' | 'neutral';
  };
  bollingerBands: {
    upper: number;
    middle: number;
    lower: number;
    position: 'above' | 'within' | 'below';
  };
  sma20: number;
  sma50: number;
  sma200: number;
  ema12: number;
  ema26: number;
  atr: number;
  stochastic: {
    k: number;
    d: number;
  };
}

export interface StockValuation {
  intrinsicValue: number;
  marginOfSafety: number;
  rating: 'undervalued' | 'fairValue' | 'overvalued';
  dcfValue?: number;
  grahamValue?: number;
  peRatio: number;
  pbRatio: number;
  psRatio: number;
  evEbitda: number;
}

export interface StockNews {
  id: string;
  title: string;
  summary: string;
  source: string;
  url: string;
  publishedAt: string;
  sentiment: 'positive' | 'neutral' | 'negative';
  symbols: string[];
}

export interface AnalystRating {
  symbol: string;
  rating: 'strong_buy' | 'buy' | 'hold' | 'sell' | 'strong_sell';
  targetPrice: number;
  numAnalysts: number;
  breakdown: {
    strongBuy: number;
    buy: number;
    hold: number;
    sell: number;
    strongSell: number;
  };
}

export interface Watchlist {
  id: string;
  name: string;
  symbols: string[];
  createdAt: string;
  updatedAt: string;
}

export interface ScreenerCriteria {
  minMarketCap?: number;
  maxMarketCap?: number;
  minPE?: number;
  maxPE?: number;
  minDividendYield?: number;
  minVolume?: number;
  sectors?: string[];
  rsiRange?: [number, number];
  priceAboveSMA?: number[];
}

export interface SectorPerformance {
  sector: string;
  change: number;
  changePercent: number;
  volume: number;
  marketCap: number;
}
